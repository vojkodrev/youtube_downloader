package main

import (
	"bytes"
	"fmt"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type ThumbnailSaver struct {
	videoDuration *VideoDuration
}

func NewThumbnailSaver(videoDuration *VideoDuration) *ThumbnailSaver {
	return &ThumbnailSaver{videoDuration: videoDuration}
}

func (ts *ThumbnailSaver) Save(videoPath, thumbnailPath string) error {
	dur, err := ts.videoDuration.Get(videoPath)
	if err != nil {
		return fmt.Errorf("probe: %w", err)
	}

	cmd := ffmpeg.Input(videoPath, ffmpeg.KwArgs{"ss": dur / 2}).
		Output(thumbnailPath, ffmpeg.KwArgs{"vframes": 1, "format": "image2"}).
		OverWriteOutput().
		Compile()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}
