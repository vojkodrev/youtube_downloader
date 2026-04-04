package main

import (
	"path/filepath"
	"strings"
)

func durationFilename(videoPath string) string {
	base := filepath.Base(videoPath)
	return strings.TrimSuffix(base, filepath.Ext(base)) + ".duration.txt"
}
