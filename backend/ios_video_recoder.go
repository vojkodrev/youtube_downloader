package main

import (
	"fmt"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// IOSVideoRecoder re-encodes a video with B-frames disabled targeting the original
// bitrate so iOS can play the file at roughly the same file size. Audio is
// stream-copied.
type IOSVideoRecoder struct{}

func NewIOSVideoRecoder() *IOSVideoRecoder {
	return &IOSVideoRecoder{}
}

func (r *IOSVideoRecoder) Recode(videoPath string, bitrate int64) error {
	ext := filepath.Ext(videoPath)
	base := videoPath[:len(videoPath)-len(ext)]
	tmpPath := base + ".temp" + ext

	cmd := ffmpeg.Input(videoPath).
		Output(tmpPath, ffmpeg.KwArgs{
			"c:v":       "libx264",
			"profile:v": "main",
			"level":     "4.0",
			"bf":        0,
			"b:v":       fmt.Sprintf("%d", bitrate),
			"vf":        "setsar=1:1",
			"c:a":       "copy",
			"movflags":  "+faststart",
		}).
		Compile()

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.Remove(videoPath); err != nil {
		return err
	}
	return os.Rename(tmpPath, videoPath)
}
