package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		videos = append(videos, Video{
			ID:       uuid.NewSHA1(uuid.NameSpaceURL, []byte(entry.Name())).String(),
			Filename: entry.Name(),
			Name:     strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())),
			Date:     info.ModTime(),
		})
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Date.After(videos[j].Date)
	})

	return videos, nil
}

func pollVideos(streamsDir string, videosMap *map[string]Video, videosMutex *sync.RWMutex) {
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
			*videosMap = m
			videosMutex.Unlock()
			log.Println("loaded", len(fetched), "videos")
		}
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	cfg := loadConfig()
	log.Println("streams dir:", cfg.StreamsDir)

	videosMap := make(map[string]Video)
	var videosMutex sync.RWMutex

	go pollVideos(cfg.StreamsDir, &videosMap, &videosMutex)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
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

	r.Run(":8080")
}
