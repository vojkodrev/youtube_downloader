package main

import "github.com/gin-gonic/gin"

type DownloadHandler struct {
	cfg   *Config
	store *VideoStore
}

func NewDownloadHandler(cfg *Config, store *VideoStore) *DownloadHandler {
	return &DownloadHandler{cfg: cfg, store: store}
}

func (h *DownloadHandler) Download(c *gin.Context) {
	id := c.Param("id")
	h.store.Mutex.RLock()
	v, ok := h.store.VideosMap[id]
	vv, versionOk := h.store.VideoVersionsMap[id]
	h.store.Mutex.RUnlock()
	if !ok && !versionOk {
		c.Status(404)
		return
	}
	var filename string
	if versionOk {
		filename = vv.Filename
	} else {
		filename = v.Filename
	}
	c.FileAttachment(h.cfg.StreamsDir+"/"+filename, filename)
}
