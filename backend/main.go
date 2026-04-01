package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StreamsDir string `yaml:"streams_dir"`
}

func loadConfig() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("could not read config.yaml:", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal("could not parse config.yaml:", err)
	}

	return cfg
}

type Video struct {
	ID       string    `json:"id"`
	Filename string    `json:"filename"`
	Name     string    `json:"name"`
	Date     time.Time `json:"date"`
}

type VideoResponse struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

func getVideos(streamsDir string) ([]Video, error) {
	entries, err := os.ReadDir(streamsDir)
	if err != nil {
		return nil, err
	}

	partRe := regexp.MustCompile(`^(.+) part\d{2}\.mp4$`)
	formatRe := regexp.MustCompile(`f\d{3}\.mp4$`)
	var videos []Video
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// skip non-mp4 files, e.g. "video.jpg", "video.mp4.duration.txt"
		if !strings.EqualFold(filepath.Ext(entry.Name()), ".mp4") {
			continue
		}
		// skip intermediate format segments, e.g. "video.f140.mp4"
		if formatRe.MatchString(entry.Name()) {
			continue
		}
		// skip in-progress downloads, e.g. "video.temp.mp4"
		if strings.HasSuffix(entry.Name(), ".temp.mp4") {
			continue
		}

		// skip part files when the merged file exists, e.g. "video part01.mp4" when "video.mp4" exists
		if m := partRe.FindStringSubmatch(entry.Name()); m != nil {
			if _, err := os.Stat(filepath.Join(streamsDir, m[1]+".mp4")); err == nil {
				continue
			}
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		v := Video{
			ID:       uuid.NewSHA1(uuid.NameSpaceURL, []byte(entry.Name())).String(),
			Filename: entry.Name(),
			Name:     name,
			Date:     info.ModTime(),
		}
		// log.Printf("video: id=%s name=%s date=%s", v.ID, v.Name, v.Date.Format("2006-01-02 15:04:05"))
		videos = append(videos, v)
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Date.After(videos[j].Date)
	})

	return videos, nil
}

func thumbnailFilename(videoName string) string {
	return videoName + ".jpg"
}

func durationFilename(videoFilename string) string {
	return videoFilename + ".duration.txt"
}

func saveThumbnail(videoPath, thumbnailPath string) error {
	dur, err := videoDuration(videoPath)
	if err != nil {
		return fmt.Errorf("probe: %w", err)
	}

	cmd := ffmpeg.Input(videoPath, ffmpeg.KwArgs{"ss": dur / 2}).
		Output(thumbnailPath, ffmpeg.KwArgs{"vframes": 1, "format": "image2"}).
		OverWriteOutput().
		Compile()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}

func videoDuration(videoPath string) (float64, error) {
	durationPath := durationFilename(videoPath)
	if data, err := os.ReadFile(durationPath); err == nil {
		return strconv.ParseFloat(string(data), 64)
	}

	probeJSON, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return 0, err
	}
	var probe struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal([]byte(probeJSON), &probe); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(probe.Format.Duration, 64)
}

func splitVideo(videoPath string) error {
	const splitDuration = 10800

	if regexp.MustCompile(`part\d{2}\.mp4$`).MatchString(videoPath) {
		return nil
	}

	dur, err := videoDuration(videoPath)
	if err != nil {
		return fmt.Errorf("probe: %w", err)
	}
	if dur <= splitDuration {
		return nil
	}

	ext := filepath.Ext(videoPath)
	base := videoPath[:len(videoPath)-len(ext)]

	cmd := ffmpeg.Input(videoPath).
		Output(base+" part%02d.mp4", ffmpeg.KwArgs{
			"c":                    "copy",
			"segment_time":         fmt.Sprintf("%d", splitDuration),
			"f":                    "segment",
			"reset_timestamps":     1,
			"segment_start_number": 1,
		}).
		OverWriteOutput().
		Compile()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	originalDir := filepath.Join(filepath.Dir(videoPath), "original")
	if err := os.MkdirAll(originalDir, 0755); err != nil {
		return err
	}
	return os.Rename(videoPath, filepath.Join(originalDir, filepath.Base(videoPath)))
}

func pollVideosWorker(streamsDir string, videos *[]Video, videosMap *map[string]Video, videosMutex *sync.RWMutex) {
	for {
		fetched, err := getVideos(streamsDir)
		if err != nil {
			log.Println("error fetching videos:", err)
		} else {
			m := make(map[string]Video, len(fetched))
			for _, v := range fetched {
				m[v.ID] = v
			}
			videosMutex.Lock()
			*videos = fetched
			*videosMap = m
			videosMutex.Unlock()
			log.Println("loaded", len(fetched), "videos")
		}
		time.Sleep(1 * time.Minute)
	}
}

func cleanupWorker(streamsDir string) {
	for {
		entries, err := os.ReadDir(streamsDir)
		if err != nil {
			log.Println("cleanup: error reading streams dir:", err)
			time.Sleep(10 * time.Minute)
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			var baseName string
			if base, ok := strings.CutSuffix(name, ".mp4.duration.txt"); ok {
				baseName = base + ".mp4"
			} else if base, ok := strings.CutSuffix(name, ".jpg"); ok {
				baseName = base + ".mp4"
			} else {
				continue
			}

			if _, err := os.Stat(filepath.Join(streamsDir, baseName)); os.IsNotExist(err) {
				path := filepath.Join(streamsDir, name)
				if err := os.Remove(path); err != nil {
					log.Println("cleanup: error removing", name, ":", err)
				} else {
					log.Println("cleanup: removed orphaned file", name)
				}
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func saveThumbnailsWorker(streamsDir string, videos *[]Video, videosMutex *sync.RWMutex) {
	for {
		videosMutex.RLock()
		current := *videos
		videosMutex.RUnlock()

		var wg sync.WaitGroup
		for _, v := range current {
			thumbPath := filepath.Join(streamsDir, thumbnailFilename(v.Name))
			// thumbnail already exists — check if it needs to be regenerated
			if tInfo, err := os.Stat(thumbPath); err == nil {
				videoPath := filepath.Join(streamsDir, v.Filename)
				newerExists := false
				// regenerate if the video file was modified after the thumbnail (e.g. file was replaced)
				if vInfo, err := os.Stat(videoPath); err == nil && vInfo.ModTime().After(tInfo.ModTime()) {
					newerExists = true
				}
				// regenerate if the duration file was updated after the thumbnail (e.g. seek position changed)
				if dInfo, err := os.Stat(filepath.Join(streamsDir, durationFilename(v.Filename))); err == nil && dInfo.ModTime().After(tInfo.ModTime()) {
					newerExists = true
				}
				if !newerExists {
					continue
				}
			}
			wg.Add(1)
			go func(v Video) {
				defer wg.Done()
				videoPath := filepath.Join(streamsDir, v.Filename)
				if err := saveThumbnail(videoPath, thumbPath); err != nil {
					log.Println("error saving thumbnail for", v.Name, ":", err)
				} else {
					log.Println("thumbnail generated for", v.Name)
				}
			}(v)
		}
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}

func saveDurationsWorker(streamsDir string, videos *[]Video, videosMutex *sync.RWMutex) {
	for {
		videosMutex.RLock()
		current := *videos
		videosMutex.RUnlock()

		var wg sync.WaitGroup
		for _, v := range current {
			durationPath := filepath.Join(streamsDir, durationFilename(v.Filename))
			// duration file already exists — skip unless the video was modified after it
			if dInfo, err := os.Stat(durationPath); err == nil {
				videoPath := filepath.Join(streamsDir, v.Filename)
				// video not newer than duration file, e.g. file was not replaced
				if vInfo, err := os.Stat(videoPath); err == nil && !vInfo.ModTime().After(dInfo.ModTime()) {
					continue
				}
			}
			wg.Add(1)
			go func(v Video) {
				defer wg.Done()
				videoPath := filepath.Join(streamsDir, v.Filename)
				duration, err := videoDuration(videoPath)
				if err != nil {
					log.Println("error getting duration for", v.Name, ":", err)
					return
				}
				if err := os.WriteFile(durationPath, []byte(strconv.FormatFloat(duration, 'f', -1, 64)), 0644); err != nil {
					log.Println("error saving duration for", v.Name, ":", err)
				} else {
					log.Println("duration saved for", v.Name)
				}
			}(v)
		}
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}

func splitVideosWorker(streamsDir string, videos *[]Video, videosMutex *sync.RWMutex) {
	for {
		videosMutex.RLock()
		current := *videos
		videosMutex.RUnlock()

		for _, v := range current {
			videoPath := filepath.Join(streamsDir, v.Filename)
			if err := splitVideo(videoPath); err != nil {
				log.Println("error splitting", v.Name, ":", err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	cfg := loadConfig()
	log.Println("streams dir:", cfg.StreamsDir)

	var videos []Video
	videosMap := make(map[string]Video)
	var videosMutex sync.RWMutex

	go pollVideosWorker(cfg.StreamsDir, &videos, &videosMap, &videosMutex)
	go func() {
		time.Sleep(5 * time.Second)
		go saveThumbnailsWorker(cfg.StreamsDir, &videos, &videosMutex)
		go saveDurationsWorker(cfg.StreamsDir, &videos, &videosMutex)
		go splitVideosWorker(cfg.StreamsDir, &videos, &videosMutex)
		go cleanupWorker(cfg.StreamsDir)
	}()

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/videos", func(c *gin.Context) {
		videosMutex.RLock()
		defer videosMutex.RUnlock()
		response := make([]VideoResponse, len(videos))
		for i, v := range videos {
			response[i] = VideoResponse{ID: v.ID, Name: v.Name, Date: v.Date}
		}
		c.JSON(200, response)
	})

	r.GET("/video/:id", func(c *gin.Context) {
		id := c.Param("id")
		videosMutex.RLock()
		v, ok := videosMap[id]
		videosMutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		c.File(cfg.StreamsDir + "/" + v.Filename)
	})

	r.GET("/download/:id", func(c *gin.Context) {
		id := c.Param("id")
		videosMutex.RLock()
		v, ok := videosMap[id]
		videosMutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		c.Header("Content-Disposition", `attachment; filename="`+v.Filename+`"`)
		c.File(cfg.StreamsDir + "/" + v.Filename)
	})

	r.GET("/duration/:id", func(c *gin.Context) {
		id := c.Param("id")
		videosMutex.RLock()
		v, ok := videosMap[id]
		videosMutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		duration, err := videoDuration(cfg.StreamsDir + "/" + v.Filename)
		if err != nil {
			c.Status(500)
			return
		}
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.JSON(200, gin.H{"duration": duration})
	})

	r.GET("/thumbnail/:id", func(c *gin.Context) {
		id := c.Param("id")
		videosMutex.RLock()
		v, ok := videosMap[id]
		videosMutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		thumbPath := filepath.Join(cfg.StreamsDir, thumbnailFilename(v.Name))
		if _, err := os.Stat(thumbPath); err != nil {
			c.Status(404)
			return
		}
		c.File(thumbPath)
	})

	r.Run(":8080")
}
