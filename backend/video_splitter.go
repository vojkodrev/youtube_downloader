package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoSplitter struct{}

func NewVideoSplitter() *VideoSplitter {
	return &VideoSplitter{}
}

func (vs *VideoSplitter) Split(videoPath string, splitDuration float64) error {
	ext := filepath.Ext(videoPath)
	base := videoPath[:len(videoPath)-len(ext)]

	cmd := ffmpeg.Input(videoPath).
		Output(base+" part%02d.mp4", ffmpeg.KwArgs{
			"c":                    "copy",
			"segment_time":         fmt.Sprintf("%g", splitDuration),
			"f":                    "segment",
			"reset_timestamps":     1,
			"segment_start_number": 1,
		}).
		OverWriteOutput().
		Compile()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	originalDir := filepath.Join(filepath.Dir(videoPath), "original")
	if err := os.MkdirAll(originalDir, 0755); err != nil {
		return err
	}
	return os.Rename(videoPath, filepath.Join(originalDir, filepath.Base(videoPath)))
}
