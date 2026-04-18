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
	videosHandler          *VideosHandler
	pingHandler            *PingHandler
	videoHandler           *VideoHandler
	downloadHandler        *DownloadHandler
	router                 *gin.Engine
}

func NewGinServer(
	cfg *Config,
	store *VideoStore,
	filenames *Filenames,
	videoDuration *VideoDuration,
	fileServer *GinSharableFileServer,
	downloadRequestService *DownloadRequestService,
	videosHandler *VideosHandler,
	pingHandler *PingHandler,
	videoHandler *VideoHandler,
	downloadHandler *DownloadHandler,
) *GinServer {
	r := gin.Default()
	return &GinServer{
		cfg:                    cfg,
		store:                  store,
		filenames:              filenames,
		videoDuration:          videoDuration,
		fileServer:             fileServer,
		downloadRequestService: downloadRequestService,
		videosHandler:          videosHandler,
		pingHandler:            pingHandler,
		videoHandler:           videoHandler,
		downloadHandler:        downloadHandler,
		router:                 r,
	}
}

func (s *GinServer) registerRoutes() {
	s.router.Use(cors.Default())

	s.router.GET("/ping", s.pingHandler.Ping)

	s.router.GET("/videos", s.videosHandler.GetVideos)

	s.router.GET("/video/:id", s.videoHandler.GetVideo)

	s.router.GET("/download/:id", s.downloadHandler.Download)

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

	s.router.POST("/delete-video", func(c *gin.Context) {
		var body struct {
			ID string `json:"id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
			return
		}

		s.store.Mutex.RLock()
		v, ok := s.store.VideosMap[body.ID]
		s.store.Mutex.RUnlock()
		if !ok {
			c.Status(404)
			return
		}

		deletedDir := filepath.Join(s.cfg.StreamsDir, "deleted")
		if err := os.MkdirAll(deletedDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create deleted folder"})
			return
		}

		// Move version files first
		for _, version := range v.Versions {
			src := filepath.Join(s.cfg.StreamsDir, version.Filename)
			if _, err := os.Stat(src); err == nil {
				os.Rename(src, filepath.Join(deletedDir, filepath.Base(version.Filename)))
			}
		}

		// Move main video file
		src := filepath.Join(s.cfg.StreamsDir, v.Filename)
		if _, err := os.Stat(src); err == nil {
			os.Rename(src, filepath.Join(deletedDir, filepath.Base(v.Filename)))
		}

		s.store.Mutex.Lock()
		delete(s.store.VideosMap, body.ID)
		for _, version := range v.Versions {
			delete(s.store.VideoVersionsMap, version.ID)
		}
		for i, sv := range s.store.Videos {
			if sv.ID == body.ID {
				s.store.Videos = append(s.store.Videos[:i], s.store.Videos[i+1:]...)
				break
			}
		}
		s.store.Mutex.Unlock()

		c.Status(http.StatusNoContent)
	})

	s.router.POST("/request-video-download", func(c *gin.Context) {
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

	s.router.POST("/request-playlist-download", func(c *gin.Context) {
		var body struct {
			URL string `json:"url" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
			return
		}

		service, playlistID, err := s.downloadRequestService.ExtractPlaylistID(body.URL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filename := fmt.Sprintf("playlist.%s.%s.download", service, playlistID)
		path := filepath.Join(s.cfg.StreamsDir, filename)

		if err := os.WriteFile(path, []byte(body.URL), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create download request"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": playlistID})
	})
}
