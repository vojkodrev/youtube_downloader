package main

import "github.com/gin-gonic/gin"

type VideoHandler struct {
	cfg        *Config
	store      *VideoStore
	fileServer *GinSharableFileServer
}

func NewVideoHandler(cfg *Config, store *VideoStore, fileServer *GinSharableFileServer) *VideoHandler {
	return &VideoHandler{cfg: cfg, store: store, fileServer: fileServer}
}

func (h *VideoHandler) GetVideo(c *gin.Context) {
	id := c.Param("id")
	h.store.Mutex.RLock()
	v, ok := h.store.VideosMap[id]
	h.store.Mutex.RUnlock()
	if !ok {
		c.Status(404)
		return
	}
	h.fileServer.Serve(c, h.cfg.StreamsDir+"/"+v.Filename, v.Filename)
}
