package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type RequestPlaylistDownloadHandler struct {
	cfg                    *Config
	downloadRequestService *DownloadRequestService
}

func NewRequestPlaylistDownloadHandler(cfg *Config, downloadRequestService *DownloadRequestService) *RequestPlaylistDownloadHandler {
	return &RequestPlaylistDownloadHandler{cfg: cfg, downloadRequestService: downloadRequestService}
}

func (h *RequestPlaylistDownloadHandler) RequestPlaylistDownload(c *gin.Context) {
	var body struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	service, playlistID, err := h.downloadRequestService.ExtractPlaylistID(body.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("playlist.%s.%s.download", service, playlistID)
	path := filepath.Join(h.cfg.StreamsDir, filename)

	if err := os.WriteFile(path, []byte(body.URL), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create download request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": playlistID})
}
