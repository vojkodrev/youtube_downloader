package main

import "github.com/gin-gonic/gin"

type VideoHandler struct {
	cfg   *Config
	store *VideoStore
}

func NewVideoHandler(cfg *Config, store *VideoStore) *VideoHandler {
	return &VideoHandler{cfg: cfg, store: store}
}

func (h *VideoHandler) GetVideo(c *gin.Context) {
	id := c.Param("id")
	h.store.Mutex.RLock()
	v, ok := h.store.VideosMap[id]
	if !ok {
		vv, vok := h.store.VideoVersionsMap[id]
		h.store.Mutex.RUnlock()
		if !vok {
			c.Status(404)
			return
		}
		c.File(h.cfg.StreamsDir + "/" + vv.Filename)
		return
	}
	h.store.Mutex.RUnlock()
	c.File(h.cfg.StreamsDir + "/" + v.Filename)
}
