package main

import (
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Video720pRecoder struct{}

func NewVideo720pRecoder() *Video720pRecoder {
	return &Video720pRecoder{}
}

func (r *Video720pRecoder) Recode(inputPath, outputPath string) error {
	tmpPath := strings.TrimSuffix(outputPath, ".mp4") + ".temp.mp4"

	cmd := ffmpeg.Input(inputPath).
		Output(tmpPath, ffmpeg.KwArgs{
			"c:v":       "libx264",
			"profile:v": "high",
			"level":     "4.0",
			"vf":        "scale=-2:min(720\\,ih)",
			"crf":       "20",
			"maxrate":   "5000000",
			"bufsize":   "10000000",
			"c:a":       "copy",
			"movflags":  "+faststart",
		}).
		OverWriteOutput().
		GlobalArgs("-loglevel", "warning", "-stats", "-stats_period", "60").
		Compile()

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return os.Rename(tmpPath, outputPath)
}
