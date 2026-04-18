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
	durationHandler        *DurationHandler
	thumbnailHandler       *ThumbnailHandler
	deleteVideoHandler     *DeleteVideoHandler
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
	durationHandler *DurationHandler,
	thumbnailHandler *ThumbnailHandler,
	deleteVideoHandler *DeleteVideoHandler,
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
		durationHandler:        durationHandler,
		thumbnailHandler:       thumbnailHandler,
		deleteVideoHandler:     deleteVideoHandler,
		router:                 r,
	}
}

func (s *GinServer) registerRoutes() {
	s.router.Use(cors.Default())

	s.router.GET("/ping", s.pingHandler.Ping)

	s.router.GET("/videos", s.videosHandler.GetVideos)

	s.router.GET("/video/:id", s.videoHandler.GetVideo)

	s.router.GET("/download/:id", s.downloadHandler.Download)

	s.router.GET("/duration/:id", s.durationHandler.GetDuration)

	s.router.GET("/thumbnail/:id", s.thumbnailHandler.GetThumbnail)

	s.router.POST("/delete-video", s.deleteVideoHandler.DeleteVideo)

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
