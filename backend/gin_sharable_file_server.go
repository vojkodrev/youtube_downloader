package main

import (
	"net/http"
	"os"
	"syscall"

	"github.com/gin-gonic/gin"
)

type GinSharableFileServer struct{}

func NewGinSharableFileServer() *GinSharableFileServer {
	return &GinSharableFileServer{}
}

func (fs *GinSharableFileServer) Serve(c *gin.Context, path string, filename string) {
	handle, err := syscall.Open(path, syscall.O_RDONLY, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_DELETE)
	if err != nil {
		c.Status(500)
		return
	}
	f := os.NewFile(uintptr(handle), path)
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		c.Status(500)
		return
	}

	http.ServeContent(c.Writer, c.Request, filename, stat.ModTime(), f)
}
