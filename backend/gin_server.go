package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type GinServer struct {
	cfg           *Config
	store         *VideoStore
	filenames     *Filenames
	videoDuration *VideoDuration
	fileServer    *GinSharableFileServer
	router        *gin.Engine
}

func NewGinServer(cfg *Config, store *VideoStore, filenames *Filenames, videoDuration *VideoDuration, fileServer *GinSharableFileServer) *GinServer {
	r := gin.Default()
	r.Use(cors.Default())

	s := &GinServer{cfg: cfg, store: store, filenames: filenames, videoDuration: videoDuration, fileServer: fileServer, router: r}
	s.registerRoutes()
	return s
}

func (s *GinServer) registerRoutes() {
	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	s.router.GET("/videos", func(c *gin.Context) {
		s.store.Mutex.RLock()
		defer s.store.Mutex.RUnlock()
		response := make([]VideoResponse, len(s.store.Videos))
		for i, v := range s.store.Videos {
			response[i] = VideoResponse{ID: v.ID, Name: v.Name, Channel: v.Channel, Date: v.Date, Status: v.Status}
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
		s.store.Mutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}
		c.Header("Content-Disposition", `attachment; filename="`+v.Filename+`"`)
		s.fileServer.Serve(c, s.cfg.StreamsDir+"/"+v.Filename, v.Filename)
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
}

func (s *GinServer) Hook(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go s.router.Run(":8080")
			return nil
		},
	})
}
