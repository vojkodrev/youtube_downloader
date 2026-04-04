package main

import (
	"encoding/json"
	"os"
	"strconv"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoDuration struct {
	filenames *Filenames
}

func NewVideoDuration(filenames *Filenames) *VideoDuration {
	return &VideoDuration{filenames: filenames}
}

func (vd *VideoDuration) Get(videoPath string) (float64, error) {
	durationPath := vd.filenames.Duration(videoPath)
	if data, err := os.ReadFile(durationPath); err == nil {
		return strconv.ParseFloat(string(data), 64)
	}

	probeJSON, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return 0, err
	}
	var probe struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if err := json.Unmarshal([]byte(probeJSON), &probe); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(probe.Format.Duration, 64)
}
