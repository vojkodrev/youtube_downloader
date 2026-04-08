package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var durationLineRegex = regexp.MustCompile(`Duration:\s+(\d+):(\d+):(\d+(?:\.\d+)?)`)

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

	probeJSON, probeErr := ffmpeg.ProbeWithTimeout(videoPath, 5*time.Second, ffmpeg.KwArgs{})
	if probeErr == nil {
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

	// ffprobe may exit with non-zero (e.g. duplicated MOOV atom on .part files)
	// but still emit duration in its output — try to parse it from the error string
	matches := durationLineRegex.FindStringSubmatch(probeErr.Error())
	if matches == nil {
		return 0, probeErr
	}
	h, _ := strconv.ParseFloat(matches[1], 64)
	m, _ := strconv.ParseFloat(matches[2], 64)
	s, _ := strconv.ParseFloat(matches[3], 64)
	total := h*3600 + m*60 + s
	if total == 0 {
		return 0, fmt.Errorf("parsed zero duration from error output: %w", probeErr)
	}
	return total, nil
}
