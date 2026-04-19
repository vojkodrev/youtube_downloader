package main

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ThumbnailHandler struct {
	cfg       *Config
	store     *VideoStore
	filenames *Filenames
}

func NewThumbnailHandler(cfg *Config, store *VideoStore, filenames *Filenames) *ThumbnailHandler {
	return &ThumbnailHandler{cfg: cfg, store: store, filenames: filenames}
}

func (h *ThumbnailHandler) GetThumbnail(c *gin.Context) {
	id := c.Param("id")
	h.store.Mutex.RLock()
	v, ok := h.store.VideosMap[id]
	h.store.Mutex.RUnlock()
	if !ok {
		c.Status(404)
		return
	}
	thumbPath := filepath.Join(h.cfg.StreamsDir, h.filenames.Thumbnail(v.Filename))
	if _, err := os.Stat(thumbPath); err != nil {
		c.Status(404)
		return
	}
	if v.Status == "Ready" {
		c.Header("Cache-Control", "public, max-age=18000")
	}
	c.File(thumbPath)
}
