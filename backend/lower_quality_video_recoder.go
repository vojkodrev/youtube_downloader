package main

import (
	"os"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type LowerQualityVideoRecoder struct{}

func NewLowerQualityVideoRecoder() *LowerQualityVideoRecoder {
	return &LowerQualityVideoRecoder{}
}

func (r *LowerQualityVideoRecoder) Recode(inputPath, outputPath string) error {
	tmpPath := strings.TrimSuffix(outputPath, ".mp4") + ".720p.temp.mp4"

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
		Compile()

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}

	if err := os.Remove(outputPath); err != nil && !os.IsNotExist(err) {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, outputPath)
}
