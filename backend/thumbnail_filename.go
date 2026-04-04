package main

import (
	"path/filepath"
	"strings"
)

func thumbnailFilename(videoPath string) string {
	base := filepath.Base(videoPath)
	return strings.TrimSuffix(base, filepath.Ext(base)) + ".jpg"
}
