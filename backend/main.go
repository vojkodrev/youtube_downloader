package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	// [download] Destination: d:\streams\CENTRIST DEMS GO TO WAR WITH HASAN, THIRD WAY CIVIL WAR, ALSO BIG PV EVENT W⧸ SENATE CANDIDATES.f299.mp4
	// [download] Destination: d:\streams\CENTRIST DEMS GO TO WAR WITH HASAN, THIRD WAY CIVIL WAR, ALSO BIG PV EVENT W⧸ SENATE CANDIDATES.f140.mp4
	entries, err := os.ReadDir(streamsDir)
	if err != nil {
		return nil, err
	}

	var videos []Video
	for _, entry := range entries {
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(entry.Name()), ".mp4") {
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

func diffVideos(videos, previousVideos []Video) (added, removed []Video) {
	prevMap := make(map[string]struct{})
	for _, v := range previousVideos {
		prevMap[v.ID] = struct{}{}
	}

	currMap := make(map[string]struct{})
	for _, v := range videos {
		currMap[v.ID] = struct{}{}
	}

	for _, v := range videos {
		if _, exists := prevMap[v.ID]; !exists {
			added = append(added, v)
		}
	}

	for _, v := range previousVideos {
		if _, exists := currMap[v.ID]; !exists {
			removed = append(removed, v)
		}
	}

	return added, removed
}

func registerVideos(r *gin.Engine, streamsDir string) {
	var previousVideos []Video

	for {
		videos, err := getVideos(streamsDir)
		if err != nil {
			log.Println("error fetching videos:", err)
		} else {
			added, removed := diffVideos(videos, previousVideos)
			previousVideos = videos

			for _, v := range removed {
				log.Println("video removed:", v.Name)
			}

			for _, v := range added {
				log.Println("video added:", v.Name)
				r.GET("/video/"+v.ID, func(c *gin.Context) {
					c.File(streamsDir + "/" + v.Filename)
				})
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	cfg := loadConfig()
	log.Println("streams dir:", cfg.StreamsDir)

	videos, err := getVideos(cfg.StreamsDir)
	if err != nil {
		log.Fatal("could not load videos:", err)
	}
	log.Println("loaded", len(videos), "videos")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
