package main

import (
	"path/filepath"
	"strings"
)

type Filenames struct{}

func NewFilenames() *Filenames {
	return &Filenames{}
}

func (f *Filenames) Thumbnail(videoPath string) string {
	base := filepath.Base(videoPath)
	return strings.TrimSuffix(base, filepath.Ext(base)) + ".jpg"
}

func (f *Filenames) Duration(videoPath string) string {
	base := filepath.Base(videoPath)
	return strings.TrimSuffix(base, filepath.Ext(base)) + ".duration.txt"
}

func (f *Filenames) IOSFix(videoPath string) string {
	base := filepath.Base(videoPath)
	return strings.TrimSuffix(base, filepath.Ext(base)) + ".ios_fix.txt"
}
