package main

import "github.com/gin-gonic/gin"

type DownloadHandler struct {
	cfg        *Config
	store      *VideoStore
	fileServer *GinSharableFileServer
}

func NewDownloadHandler(cfg *Config, store *VideoStore, fileServer *GinSharableFileServer) *DownloadHandler {
	return &DownloadHandler{cfg: cfg, store: store, fileServer: fileServer}
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
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	h.fileServer.Serve(c, h.cfg.StreamsDir+"/"+filename, filename)
}
