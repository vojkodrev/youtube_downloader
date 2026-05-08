package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type GinServer struct {
	cfg                    *Config
	store                  *VideoStore
	filenames              *Filenames
	videoDuration          *VideoDuration
	downloadRequestService *DownloadRequestService
	authHandler            *AuthHandler
	authMiddleware         *AuthMiddleware
	videosHandler          *VideosHandler
	pingHandler            *PingHandler
	videoHandler           *VideoHandler
	downloadHandler        *DownloadHandler
	durationHandler        *DurationHandler
	thumbnailHandler       *ThumbnailHandler
	deleteVideoHandler              *DeleteVideoHandler
	requestVideoDownloadHandler    *RequestVideoDownloadHandler
	requestPlaylistDownloadHandler *RequestPlaylistDownloadHandler
	router                         *gin.Engine
}

func NewGinServer(
	cfg *Config,
	store *VideoStore,
	filenames *Filenames,
	videoDuration *VideoDuration,
	downloadRequestService *DownloadRequestService,
	authHandler *AuthHandler,
	authMiddleware *AuthMiddleware,
	videosHandler *VideosHandler,
	pingHandler *PingHandler,
	videoHandler *VideoHandler,
	downloadHandler *DownloadHandler,
	durationHandler *DurationHandler,
	thumbnailHandler *ThumbnailHandler,
	deleteVideoHandler *DeleteVideoHandler,
	requestVideoDownloadHandler *RequestVideoDownloadHandler,
	requestPlaylistDownloadHandler *RequestPlaylistDownloadHandler,
) *GinServer {
	r := gin.Default()
	return &GinServer{
		cfg:                            cfg,
		store:                          store,
		filenames:                      filenames,
		videoDuration:                  videoDuration,
		downloadRequestService:         downloadRequestService,
		authHandler:                    authHandler,
		authMiddleware:                 authMiddleware,
		videosHandler:                  videosHandler,
		pingHandler:                    pingHandler,
		videoHandler:                   videoHandler,
		downloadHandler:                downloadHandler,
		durationHandler:                durationHandler,
		thumbnailHandler:               thumbnailHandler,
		deleteVideoHandler:             deleteVideoHandler,
		requestVideoDownloadHandler:    requestVideoDownloadHandler,
		requestPlaylistDownloadHandler: requestPlaylistDownloadHandler,
		router:                         r,
	}
}

func (s *GinServer) registerRoutes() {
	s.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{s.cfg.FrontendURL},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	s.router.GET("/ping", s.pingHandler.Ping)
	s.router.GET("/auth/login", s.authHandler.Login)
	s.router.GET("/auth/callback", s.authHandler.Callback)

	protected := s.router.Group("/", s.authMiddleware.Handle)

	protected.GET("/videos", s.videosHandler.GetVideos)

	protected.GET("/video/:id", s.videoHandler.GetVideo)

	protected.GET("/download/:id", s.downloadHandler.Download)

	protected.GET("/duration/:id", s.durationHandler.GetDuration)

	protected.GET("/thumbnail/:id", s.thumbnailHandler.GetThumbnail)

	protected.POST("/delete-video", s.deleteVideoHandler.DeleteVideo)

	protected.POST("/request-video-download", s.requestVideoDownloadHandler.RequestVideoDownload)

	protected.POST("/request-playlist-download", s.requestPlaylistDownloadHandler.RequestPlaylistDownload)
}
