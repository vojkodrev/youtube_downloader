package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type RequestVideoDownloadHandler struct {
	cfg                    *Config
	downloadRequestService *DownloadRequestService
}

func NewRequestVideoDownloadHandler(cfg *Config, downloadRequestService *DownloadRequestService) *RequestVideoDownloadHandler {
	return &RequestVideoDownloadHandler{cfg: cfg, downloadRequestService: downloadRequestService}
}

func (h *RequestVideoDownloadHandler) RequestVideoDownload(c *gin.Context) {
	var body struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	if !h.downloadRequestService.IsSupportedURL(body.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Url must be from a supported service (YouTube, Twitch)"})
		return
	}

	service, videoID, err := h.downloadRequestService.ExtractVideoID(body.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := fmt.Sprintf("video.%s.%s.download", service, videoID)
	path := filepath.Join(h.cfg.StreamsDir, filename)

	if err := os.WriteFile(path, []byte(body.URL), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create download request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": videoID})
}
