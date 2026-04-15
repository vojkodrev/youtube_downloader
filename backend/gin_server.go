package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	cfg                    *Config
	store                  *VideoStore
	filenames              *Filenames
	videoDuration          *VideoDuration
	fileServer             *GinSharableFileServer
	downloadRequestService *DownloadRequestService
	router                 *gin.Engine
}

func NewGinServer(cfg *Config, store *VideoStore, filenames *Filenames, videoDuration *VideoDuration, fileServer *GinSharableFileServer, downloadRequestService *DownloadRequestService) *GinServer {
	r := gin.Default()
	return &GinServer{cfg: cfg, store: store, filenames: filenames, videoDuration: videoDuration, fileServer: fileServer, downloadRequestService: downloadRequestService, router: r}
}

func (s *GinServer) registerRoutes() {
	s.router.Use(cors.Default())

	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	s.router.GET("/videos", func(c *gin.Context) {
		s.store.Mutex.RLock()
		defer s.store.Mutex.RUnlock()
		response := make([]VideoResponse, len(s.store.Videos))
		for i, v := range s.store.Videos {
			response[i] = VideoResponse{
				ID:       v.ID,
				Name:     v.Name,
				Channel:  v.Channel,
				Date:     v.Date,
				Status:   v.Status,
				Versions: v.Versions,
			}
		}
		c.JSON(200, response)
	})

	s.router.GET("/video/:id", func(c *gin.Context) {
		id := c.Param("id")
		s.store.Mutex.RLock()
		v, ok := s.store.VideosMap[id]
		s.store.Mutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		s.fileServer.Serve(c, s.cfg.StreamsDir+"/"+v.Filename, v.Filename)
	})

	s.router.GET("/download/:id", func(c *gin.Context) {
		id := c.Param("id")
		s.store.Mutex.RLock()
		v, ok := s.store.VideosMap[id]
		vv, versionOk := s.store.VideoVersionsMap[id]
		s.store.Mutex.RUnlock()
		if !ok && !versionOk {
			c.Status(404)
			return
		}
		var filename string
		if versionOk {
			filename = vv.Filename
		} else {
			filename = v.Filename
		}
		c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
		s.fileServer.Serve(c, s.cfg.StreamsDir+"/"+filename, filename)
	})

	s.router.GET("/duration/:id", func(c *gin.Context) {
		id := c.Param("id")
		s.store.Mutex.RLock()
		v, ok := s.store.VideosMap[id]
		s.store.Mutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		duration, err := s.videoDuration.Get(s.cfg.StreamsDir + "/" + v.Filename)
		if err != nil {
			c.Status(500)
			return
		}
		if v.Status == "Ready" {
			c.Header("Cache-Control", "public, max-age=18000")
		}
		c.JSON(200, gin.H{"duration": duration})
	})

	s.router.GET("/thumbnail/:id", func(c *gin.Context) {
		id := c.Param("id")
		s.store.Mutex.RLock()
		v, ok := s.store.VideosMap[id]
		s.store.Mutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		thumbPath := filepath.Join(s.cfg.StreamsDir, s.filenames.Thumbnail(v.Filename))
		if _, err := os.Stat(thumbPath); err != nil {
			c.Status(404)
			return
		}
		if v.Status == "Ready" {
			c.Header("Cache-Control", "public, max-age=18000")
		}
		c.File(thumbPath)
	})

	s.router.POST("/request-download", func(c *gin.Context) {
		var body struct {
			URL string `json:"url" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}

		if !s.downloadRequestService.IsSupportedURL(body.URL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Url must be from a supported service (YouTube, Twitch)"})
			return
		}

		service, videoID, err := s.downloadRequestService.ExtractVideoID(body.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filename := fmt.Sprintf("video.%s.%s.download", service, videoID)
		path := filepath.Join(s.cfg.StreamsDir, filename)

		if err := os.WriteFile(path, []byte(body.URL), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create download request"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": videoID})
	})
}
