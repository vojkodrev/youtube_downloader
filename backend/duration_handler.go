package main

import "github.com/gin-gonic/gin"

type DurationHandler struct {
	cfg           *Config
	store         *VideoStore
	videoDuration *VideoDuration
}

func NewDurationHandler(cfg *Config, store *VideoStore, videoDuration *VideoDuration) *DurationHandler {
	return &DurationHandler{cfg: cfg, store: store, videoDuration: videoDuration}
}

func (h *DurationHandler) GetDuration(c *gin.Context) {
	id := c.Param("id")
	h.store.Mutex.RLock()
	v, ok := h.store.VideosMap[id]
	h.store.Mutex.RUnlock()
	if !ok {
		c.Status(404)
		return
	}
	duration, err := h.videoDuration.Get(h.cfg.StreamsDir + "/" + v.Filename)
	if err != nil {
		c.Status(500)
		return
	}
	if v.Status == "Ready" {
		c.Header("Cache-Control", "public, max-age=18000")
	}
	c.JSON(200, gin.H{"duration": duration})
}
