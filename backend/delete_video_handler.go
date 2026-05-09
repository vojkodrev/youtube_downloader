package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type DeleteVideoHandler struct {
	cfg   *Config
	store *VideoStore
}

func NewDeleteVideoHandler(cfg *Config, store *VideoStore) *DeleteVideoHandler {
	return &DeleteVideoHandler{cfg: cfg, store: store}
}

func (h *DeleteVideoHandler) DeleteVideo(c *gin.Context) {
	var body struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	h.store.Mutex.RLock()
	v, ok := h.store.VideosMap[body.ID]
	h.store.Mutex.RUnlock()
	if !ok {
		c.Status(404)
		return
	}

	for _, version := range v.Versions {
		src := filepath.Join(h.cfg.StreamsDir, version.Filename)
		if _, err := os.Stat(src); err == nil {
			_ = os.Remove(src)
		}
	}

	src := filepath.Join(h.cfg.StreamsDir, v.Filename)
	if _, err := os.Stat(src); err == nil {
		_ = os.Remove(src)
	}

	h.store.Mutex.Lock()
	delete(h.store.VideosMap, body.ID)
	for _, version := range v.Versions {
		delete(h.store.VideoVersionsMap, version.ID)
	}
	for i, sv := range h.store.Videos {
		if sv.ID == body.ID {
			h.store.Videos = append(h.store.Videos[:i], h.store.Videos[i+1:]...)
			break
		}
	}
	h.store.Mutex.Unlock()

	c.Status(http.StatusNoContent)
}
