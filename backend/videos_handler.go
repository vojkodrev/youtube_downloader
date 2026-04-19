package main

import "github.com/gin-gonic/gin"

type VideosHandler struct {
	store *VideoStore
}

func NewVideosHandler(store *VideoStore) *VideosHandler {
	return &VideosHandler{store: store}
}

func (h *VideosHandler) GetVideos(c *gin.Context) {
	h.store.Mutex.RLock()
	defer h.store.Mutex.RUnlock()
	response := make([]VideoResponse, len(h.store.Videos))
	for i, v := range h.store.Videos {
		response[i] = VideoResponse{
			ID:       v.ID,
			Name:     v.Name,
			Channel:  v.Channel,
			Date:     v.Date,
			Status:   v.Status,
			Versions: v.Versions,
		}
	}
	c.JSON(200, response)
}
