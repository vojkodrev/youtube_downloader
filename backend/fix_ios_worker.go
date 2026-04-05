package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FixIOSWorker struct {
	cfg       *Config
	store     *VideoStore
	filenames *Filenames
}

func NewFixIOSWorker(cfg *Config, store *VideoStore, filenames *Filenames) *FixIOSWorker {
	return &FixIOSWorker{cfg: cfg, store: store, filenames: filenames}
}

func (fw *FixIOSWorker) Start() {
	for {
		fw.store.Mutex.RLock()
		current := fw.store.Videos
		fw.store.Mutex.RUnlock()

		for _, v := range current {
			if v.Status != "Ready" {
				continue
			}
			videoPath := filepath.Join(fw.cfg.StreamsDir, v.Filename)
			fixPath := filepath.Join(fw.cfg.StreamsDir, fw.filenames.IOSFix(v.Filename))
			// skip if marker exists and video has not been replaced since
			if fixInfo, err := os.Stat(fixPath); err == nil {
				if vInfo, err := os.Stat(videoPath); err != nil || !vInfo.ModTime().After(fixInfo.ModTime()) {
					continue
				}
			}
			needsFix, err := fw.needsFix(videoPath)
			if err != nil {
				log.Println("error probing", v.Filename, ":", err)
				continue
			}
			if !needsFix {
				if err := os.WriteFile(fixPath, []byte("ok"), 0644); err != nil {
					log.Println("error writing ios fix marker for", v.Filename, ":", err)
				}
				continue
			}
			log.Println("fixing ios compatibility for", v.Filename)
			if err := fw.fix(videoPath); err != nil {
				log.Println("error fixing", v.Filename, ":", err)
				continue
			}
			if err := os.WriteFile(fixPath, []byte("fixed sar"), 0644); err != nil {
				log.Println("error writing ios fix marker for", v.Filename, ":", err)
			}
			log.Println("ios fix done for", v.Filename)
		}
		time.Sleep(1 * time.Minute)
	}
}

// needsFix probes the video and returns true when the stream has B-frames or
// a missing SAR, either of which causes iOS to refuse playback.
// A working file has has_b_frames=0 and a valid sample_aspect_ratio like "1:1".
func (fw *FixIOSWorker) needsFix(videoPath string) (bool, error) {
	probeJSON, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return false, err
	}
	var probe struct {
		Streams []struct {
			CodecType  string `json:"codec_type"`
			SAR        string `json:"sample_aspect_ratio"`
			HasBFrames int    `json:"has_b_frames"`
		} `json:"streams"`
	}
	if err := json.Unmarshal([]byte(probeJSON), &probe); err != nil {
		return false, err
	}
	for _, s := range probe.Streams {
		if s.CodecType != "video" {
			continue
		}
		log.Printf("sample_aspect_ratio=%q has_b_frames=%d for %s", s.SAR, s.HasBFrames, videoPath)
		missingSAR := s.SAR == "" || s.SAR == "N/A" || s.SAR == "0:1"
		return missingSAR || s.HasBFrames > 0, nil
	}
	return false, fmt.Errorf("no video stream found in %s", videoPath)
}

// fix re-encodes the video with B-frames disabled and an explicit SAR of 1:1
// so iOS can play the file. Audio is stream-copied unchanged.
func (fw *FixIOSWorker) fix(videoPath string) error {
	ext := filepath.Ext(videoPath)
	base := videoPath[:len(videoPath)-len(ext)]
	tmpPath := base + ".temp" + ext

	cmd := ffmpeg.Input(videoPath).
		Output(tmpPath, ffmpeg.KwArgs{
			"c:v":      "libx264",
			"profile:v": "main",
			"level":    "4.0",
			"bf":       0,
			"c:a":      "copy",
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
