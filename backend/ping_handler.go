package main

import "github.com/gin-gonic/gin"

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
