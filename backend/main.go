package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	Channel  string    `json:"channel"`
	Date     time.Time `json:"date"`
	Status   string    `json:"status"`
}

type VideoResponse struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Channel string    `json:"channel"`
	Date    time.Time `json:"date"`
	Status  string    `json:"status"`
}

func getVideos(streamsDir string) ([]Video, error) {
	entries, err := os.ReadDir(streamsDir)
	if err != nil {
		return nil, err
	}

	partRe := regexp.MustCompile(`^(.+) part\d{2}\.mp4$`)
	formatRe := regexp.MustCompile(`f\d{3}\.mp4$`)
	ytdlRe := regexp.MustCompile(`\.f\d{3}\.[^.]+\.ytdl$`)
	channelRe := regexp.MustCompile(`^\[([^\]]+)\] ?`)

	// pre-scan: for each ytdl base name, find the largest file
	largestYtdl := map[string]string{} // base -> filename with largest size
	largestYtdlSize := map[string]int64{}
	for _, entry := range entries {
		if strings.ToLower(filepath.Ext(entry.Name())) != ".ytdl" {
			continue
		}
		base := ytdlRe.ReplaceAllString(entry.Name(), "")
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.Size() > largestYtdlSize[base] {
			largestYtdlSize[base] = info.Size()
			largestYtdl[base] = entry.Name()
		}
	}

	var videos []Video
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// skip non-mp4/ytdl files, e.g. "video.jpg", "video.mp4.duration.txt"
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".mp4" && ext != ".ytdl" {
			continue
		}
		status := "Ready"
		if ext == ".ytdl" {
			// only include the largest ytdl file per base (skip smaller format segments)
			base := ytdlRe.ReplaceAllString(entry.Name(), "")
			if largestYtdl[base] != entry.Name() {
				continue
			}
			status = "Downloading"
		}
		// skip intermediate format segments, e.g. "video.f140.mp4"
		if formatRe.MatchString(entry.Name()) {
			continue
		}
		// in-progress yt-dlp download
		if strings.HasSuffix(entry.Name(), ".temp.mp4") {
			status = "Downloading"
		}
		if m := partRe.FindStringSubmatch(entry.Name()); m != nil {
			// this is a partXX file — skip it if the source file still exists (splitting in progress)
			if _, err := os.Stat(filepath.Join(streamsDir, m[1]+".mp4")); err == nil {
				continue
			}
		} else {
			// this is a plain mp4 — mark as Processing if any partXX files exist (splitting in progress)
			base := strings.TrimSuffix(entry.Name(), ".mp4")
			pattern := filepath.Join(streamsDir, base+" part*.mp4")
			if matches, _ := filepath.Glob(pattern); len(matches) > 0 {
				status = "Processing"
			}
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		var channel string
		if m := channelRe.FindStringSubmatch(name); m != nil {
			channel = m[1]
			name = name[len(m[0]):]
		}
		v := Video{
			ID:       uuid.NewSHA1(uuid.NameSpaceURL, []byte(entry.Name())).String(),
			Filename: entry.Name(),
			Name:     name,
			Channel:  channel,
			Date:     info.ModTime(),
			Status:   status,
		}
		// log.Printf("video: id=%s name=%s date=%s", v.ID, v.Name, v.Date.Format("2006-01-02 15:04:05"))
		videos = append(videos, v)
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Date.After(videos[j].Date)
	})

	return videos, nil
}

func thumbnailFilename(videoPath string) string {
	base := filepath.Base(videoPath)
	return strings.TrimSuffix(base, filepath.Ext(base)) + ".jpg"
}

func durationFilename(videoPath string) string {
	base := filepath.Base(videoPath)
	return strings.TrimSuffix(base, filepath.Ext(base)) + ".duration.txt"
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
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			var base string
			if b, ok := strings.CutSuffix(name, ".duration.txt"); ok {
				base = b
			} else if b, ok := strings.CutSuffix(name, ".jpg"); ok {
				base = b
			} else {
				continue
			}

			exists := false
			for _, ext := range []string{".mp4", ".ytdl"} {
				if _, err := os.Stat(filepath.Join(streamsDir, base+ext)); err == nil {
					exists = true
					break
				}
			}
			if !exists {
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
			thumbPath := filepath.Join(streamsDir, thumbnailFilename(v.Filename))
			videoPath := filepath.Join(streamsDir, v.Filename)
			durationPath := filepath.Join(streamsDir, durationFilename(v.Filename))
			// thumbnail already exists — check if it needs to be regenerated
			if tInfo, err := os.Stat(thumbPath); err == nil {
				newerExists := false
				// regenerate if the video file was modified after the thumbnail (e.g. file was replaced)
				if vInfo, err := os.Stat(videoPath); err == nil && vInfo.ModTime().After(tInfo.ModTime()) {
					newerExists = true
				}
				// regenerate if the duration file was updated after the thumbnail (e.g. seek position changed)
				if dInfo, err := os.Stat(durationPath); err == nil && dInfo.ModTime().After(tInfo.ModTime()) {
					newerExists = true
				}
				if !newerExists {
					continue
				}
			}
			wg.Go(func() {
				if err := saveThumbnail(videoPath, thumbPath); err != nil {
					log.Println("error saving thumbnail for", v.Filename, ":", err)
				} else {
					log.Println("thumbnail generated for", v.Filename)
				}
			})
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
			videoPath := filepath.Join(streamsDir, v.Filename)
			// duration file already exists — skip unless the video was modified after it
			if dInfo, err := os.Stat(durationPath); err == nil {
				// video not newer than duration file, e.g. file was not replaced
				if vInfo, err := os.Stat(videoPath); err == nil && !vInfo.ModTime().After(dInfo.ModTime()) {
					continue
				}
			}
			wg.Go(func() {
				duration, err := videoDuration(videoPath)
				if err != nil {
					log.Println("error getting duration for", v.Filename, ":", err)
					return
				}
				if err := os.WriteFile(durationPath, []byte(strconv.FormatFloat(duration, 'f', -1, 64)), 0644); err != nil {
					log.Println("error saving duration for", v.Filename, ":", err)
				} else {
					log.Println("duration saved for", v.Filename)
				}
			})
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

		partRe := regexp.MustCompile(` part\d{2}\.mp4$`)
		for _, v := range current {
			if v.Status != "Ready" || partRe.MatchString(v.Filename) {
				continue
			}
			videoPath := filepath.Join(streamsDir, v.Filename)
			if err := splitVideo(videoPath); err != nil {
				log.Println("error splitting", v.Filename, ":", err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func serveFileSharable(c *gin.Context, path string, filename string) {
	handle, err := syscall.Open(path, syscall.O_RDONLY, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_DELETE)
	if err != nil {
		c.Status(500)
		return
	}
	f := os.NewFile(uintptr(handle), path)
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		c.Status(500)
		return
	}

	http.ServeContent(c.Writer, c.Request, filename, stat.ModTime(), f)
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
			response[i] = VideoResponse{ID: v.ID, Name: v.Name, Channel: v.Channel, Date: v.Date, Status: v.Status}
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
		serveFileSharable(c, cfg.StreamsDir+"/"+v.Filename, v.Filename)
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
		serveFileSharable(c, cfg.StreamsDir+"/"+v.Filename, v.Filename)
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
		c.Header("Cache-Control", "public, max-age=18000")
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
		thumbPath := filepath.Join(cfg.StreamsDir, thumbnailFilename(v.Filename))
		if _, err := os.Stat(thumbPath); err != nil {
			c.Status(404)
			return
		}
		c.Header("Cache-Control", "public, max-age=18000")
		c.File(thumbPath)
	})

	r.Run(":8080")
}
