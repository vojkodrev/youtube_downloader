package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
		if v.Status == "Ready" {
			c.Header("Cache-Control", "public, max-age=18000")
		}
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
		if v.Status == "Ready" {
			c.Header("Cache-Control", "public, max-age=18000")
		}
		c.File(thumbPath)
	})

	r.Run(":8080")
}
