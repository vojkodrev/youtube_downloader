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
	fileServer             *GinSharableFileServer
	downloadRequestService *DownloadRequestService
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
	fileServer *GinSharableFileServer,
	downloadRequestService *DownloadRequestService,
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
		fileServer:                     fileServer,
		downloadRequestService:         downloadRequestService,
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
	s.router.Use(cors.Default())

	s.router.GET("/ping", s.pingHandler.Ping)

	s.router.GET("/videos", s.videosHandler.GetVideos)

	s.router.GET("/video/:id", s.videoHandler.GetVideo)

	s.router.GET("/download/:id", s.downloadHandler.Download)

	s.router.GET("/duration/:id", s.durationHandler.GetDuration)

	s.router.GET("/thumbnail/:id", s.thumbnailHandler.GetThumbnail)

	s.router.POST("/delete-video", s.deleteVideoHandler.DeleteVideo)

	s.router.POST("/request-video-download", s.requestVideoDownloadHandler.RequestVideoDownload)

	s.router.POST("/request-playlist-download", s.requestPlaylistDownloadHandler.RequestPlaylistDownload)
}
