package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
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
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Name      string    `json:"name"`
	Date      time.Time `json:"date"`
	Thumbnail string    `json:"thumbnail_path"`
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

	var videos []Video
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".mp4") {
			continue
		}
		if strings.HasSuffix(entry.Name(), "f299.mp4") || strings.HasSuffix(entry.Name(), "f140.mp4") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		thumbFilename := thumbnailFilename(name)
		thumb := ""
		if _, err := os.Stat(filepath.Join(streamsDir, thumbFilename)); err == nil {
			thumb = thumbFilename
		}

		v := Video{
			ID:        uuid.NewSHA1(uuid.NameSpaceURL, []byte(entry.Name())).String(),
			Filename:  entry.Name(),
			Name:      name,
			Date:      info.ModTime(),
			Thumbnail: thumb,
		}
		log.Printf("video: id=%s name=%s thumbnail=%s date=%s", v.ID, v.Name, v.Thumbnail, v.Date.Format("2006-01-02 15:04:05"))
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

func saveThumbnail(videoPath, thumbnailPath string) error {
	if _, err := os.Stat(thumbnailPath); err == nil {
		return nil
	}

	return ffmpeg.Input(videoPath, ffmpeg.KwArgs{"ss": 10}).
		Output(thumbnailPath, ffmpeg.KwArgs{"vframes": 1, "format": "image2"}).
		OverWriteOutput().
		Run()
}

func pollVideos(streamsDir string, videos *[]Video, videosMap *map[string]Video, videosMutex *sync.RWMutex) {
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

func pollThumbnails(streamsDir string, videos *[]Video, videosMutex *sync.RWMutex) {
	for {
		videosMutex.RLock()
		current := *videos
		videosMutex.RUnlock()

		for _, v := range current {
			if v.Thumbnail != "" {
				continue
			}
			videoPath := filepath.Join(streamsDir, v.Filename)
			thumbPath := filepath.Join(streamsDir, thumbnailFilename(v.Name))
			if err := saveThumbnail(videoPath, thumbPath); err != nil {
				log.Println("error saving thumbnail for", v.Name, ":", err)
			} else {
				log.Println("thumbnail generated for", v.Name)
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

	go pollVideos(cfg.StreamsDir, &videos, &videosMap, &videosMutex)
	go pollThumbnails(cfg.StreamsDir, &videos, &videosMutex)

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

	r.GET("/thumbnail/:id", func(c *gin.Context) {
		id := c.Param("id")
		videosMutex.RLock()
		v, ok := videosMap[id]
		videosMutex.RUnlock()
		if !ok || v.Thumbnail == "" {
			c.Status(404)
			return
		}
		c.File(filepath.Join(cfg.StreamsDir, v.Thumbnail))
	})

	r.Run(":8080")
}
