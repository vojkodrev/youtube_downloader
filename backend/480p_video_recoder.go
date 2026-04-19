package main

import (
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Video480pRecoder struct{}

func NewVideo480pRecoder() *Video480pRecoder {
	return &Video480pRecoder{}
}

func (r *Video480pRecoder) Recode(inputPath, outputPath string) error {
	tmpPath := strings.TrimSuffix(outputPath, ".mp4") + ".temp.mp4"

	cmd := ffmpeg.Input(inputPath).
		Output(tmpPath, ffmpeg.KwArgs{
			"c:v":       "libx264",
			"profile:v": "high",
			"level":     "3.1",
			"vf":        "scale=-2:min(480\\,ih)",
			"crf":       "22",
			"maxrate":   "2500000",
			"bufsize":   "5000000",
			"c:a":       "copy",
			"movflags":  "+faststart",
		}).
		GlobalArgs("-loglevel", "warning", "-stats", "-stats_period", "60").
		Compile()

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return os.Rename(tmpPath, outputPath)
}
