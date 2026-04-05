package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FixFrameRateWorker struct {
	cfg       *Config
	store     *VideoStore
	filenames *Filenames
}

func NewFixFrameRateWorker(cfg *Config, store *VideoStore, filenames *Filenames) *FixFrameRateWorker {
	return &FixFrameRateWorker{cfg: cfg, store: store, filenames: filenames}
}

func (fw *FixFrameRateWorker) Start() {
	for {
		fw.store.Mutex.RLock()
		current := fw.store.Videos
		fw.store.Mutex.RUnlock()

		for _, v := range current {
			if v.Status != "Ready" {
				continue
			}
			videoPath := filepath.Join(fw.cfg.StreamsDir, v.Filename)
			fixPath := filepath.Join(fw.cfg.StreamsDir, fw.filenames.FrameRateFix(v.Filename))
			// skip if marker exists and video has not been replaced since
			if fixInfo, err := os.Stat(fixPath); err == nil {
				if vInfo, err := os.Stat(videoPath); err != nil || !vInfo.ModTime().After(fixInfo.ModTime()) {
					continue
				}
			}
			targetFPS, needsFix, err := fw.checkFrameRate(videoPath)
			if err != nil {
				log.Println("error probing frame rate for", v.Filename, ":", err)
				continue
			}
			if !needsFix {
				if err := os.WriteFile(fixPath, []byte("ok"), 0644); err != nil {
					log.Println("error writing frame rate fix marker for", v.Filename, ":", err)
				}
				continue
			}
			log.Printf("fixing frame rate of %s to %d fps", v.Filename, targetFPS)
			if err := fw.fixFrameRate(videoPath, targetFPS); err != nil {
				log.Println("error fixing frame rate for", v.Filename, ":", err)
				continue
			}
			if err := os.WriteFile(fixPath, []byte(fmt.Sprintf("fixed to %d fps", targetFPS)), 0644); err != nil {
				log.Println("error writing frame rate fix marker for", v.Filename, ":", err)
			}
			log.Println("frame rate fixed for", v.Filename)
		}
		time.Sleep(1 * time.Minute)
	}
}

// checkFrameRate probes the video and returns (targetFPS, needsFix, error).
// targetFPS is 30 for standard content and 60 for high-frame-rate content.
// needsFix is true when the actual frame rate is not a whole-number multiple of the target.
func (fw *FixFrameRateWorker) checkFrameRate(videoPath string) (int, bool, error) {
	probeJSON, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return 0, false, err
	}
	var probe struct {
		Streams []struct {
			CodecType string `json:"codec_type"`
			AvgFPS    string `json:"avg_frame_rate"`
			RFps      string `json:"r_frame_rate"`
		} `json:"streams"`
	}
	if err := json.Unmarshal([]byte(probeJSON), &probe); err != nil {
		return 0, false, err
	}
	for _, s := range probe.Streams {
		if s.CodecType != "video" {
			continue
		}
		// avg_frame_rate is the average fps over the whole file, e.g. "431808000/14392231" for a VFR stream or "30/1" for CFR
		fps, err := parseFraction(s.AvgFPS)
		if err != nil || fps == 0 {
			// r_frame_rate is the container's declared fps, e.g. "30/1" — fall back when avg is missing or zero
			fps, err = parseFraction(s.RFps)
			if err != nil || fps == 0 {
				return 0, false, fmt.Errorf("could not determine fps from %q / %q", s.AvgFPS, s.RFps)
			}
		}
		target := 30
		if fps > 35 {
			target = 60
		}
		// consider it already correct when fps is within 0.1% of the target
		if math.Abs(fps-float64(target))/float64(target) < 0.001 {
			return target, false, nil
		}
		return target, true, nil
	}
	return 0, false, fmt.Errorf("no video stream found in %s", videoPath)
}

func (fw *FixFrameRateWorker) fixFrameRate(videoPath string, targetFPS int) error {
	ext := filepath.Ext(videoPath)
	base := videoPath[:len(videoPath)-len(ext)]
	tmpPath := base + ".temp" + ext

	cmd := ffmpeg.Input(videoPath).
		Output(tmpPath, ffmpeg.KwArgs{
			"vf":   fmt.Sprintf("fps=%d", targetFPS),
			"c:v":  "libx264",
			"c:a":  "copy",
			"movflags": "+faststart",
		}).
		OverWriteOutput().
		Compile()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	originalDir := filepath.Join(filepath.Dir(videoPath), "original")
	if err := os.MkdirAll(originalDir, 0755); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(videoPath, filepath.Join(originalDir, filepath.Base(videoPath))); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, videoPath)
}

func parseFraction(s string) (float64, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return strconv.ParseFloat(s, 64)
	}
	num, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}
	den, err := strconv.ParseFloat(parts[1], 64)
	if err != nil || den == 0 {
		return 0, fmt.Errorf("invalid denominator in %q", s)
	}
	return num / den, nil
}
