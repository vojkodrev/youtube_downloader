package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoDuration struct {
	cfg       *Config
	filenames *Filenames
}

func NewVideoDuration(cfg *Config, filenames *Filenames) *VideoDuration {
	return &VideoDuration{cfg: cfg, filenames: filenames}
}

func (vd *VideoDuration) Get(videoPath string) (float64, error) {
	return vd.get(videoPath, false)
}

func (vd *VideoDuration) GetForce(videoPath string) (float64, error) {
	return vd.get(videoPath, true)
}

func (vd *VideoDuration) get(videoPath string, force bool) (float64, error) {
	durationPath := filepath.Join(vd.cfg.StreamsDir, vd.filenames.Duration(videoPath))
	if !force {
		if data, err := os.ReadFile(durationPath); err == nil {
			return strconv.ParseFloat(string(data), 64)
		}
	}

	probeJSON, err := ffmpeg.ProbeWithTimeout(videoPath, 45*time.Second, ffmpeg.KwArgs{})
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
