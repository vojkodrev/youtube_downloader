package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FixIOSWorker struct {
	cfg           *Config
	store         *VideoStore
	filenames     *Filenames
	videoDuration *VideoDuration
}

func NewFixIOSWorker(cfg *Config, store *VideoStore, filenames *Filenames, videoDuration *VideoDuration) *FixIOSWorker {
	return &FixIOSWorker{cfg: cfg, store: store, filenames: filenames, videoDuration: videoDuration}
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
			dur, err := fw.videoDuration.Get(videoPath)
			if err != nil || dur > fw.cfg.SplitDuration {
				continue
			}
			fixPath := filepath.Join(fw.cfg.StreamsDir, fw.filenames.IOSFix(v.Filename))
			// skip if marker exists and video has not been replaced since
			if fixInfo, err := os.Stat(fixPath); err == nil {
				if vInfo, err := os.Stat(videoPath); err != nil || !vInfo.ModTime().After(fixInfo.ModTime()) {
					continue
				}
			}
			needsFix, bitrate, err := fw.needsFix(videoPath)
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
			if err := fw.fix(videoPath, bitrate); err != nil {
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

// needsFix probes the video and returns (needsFix, videoBitrate, error).
// needsFix is true when the stream has B-frames or a missing SAR, either of
// which causes iOS to refuse playback.
// A working file has has_b_frames=0 and a valid sample_aspect_ratio like "1:1".
func (fw *FixIOSWorker) needsFix(videoPath string) (bool, int64, error) {
	probeJSON, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return false, 0, err
	}
	var probe struct {
		Streams []struct {
			CodecType  string `json:"codec_type"`
			SAR        string `json:"sample_aspect_ratio"`
			HasBFrames int    `json:"has_b_frames"`
			BitRate    string `json:"bit_rate"`
		} `json:"streams"`
	}
	if err := json.Unmarshal([]byte(probeJSON), &probe); err != nil {
		return false, 0, err
	}
	for _, s := range probe.Streams {
		if s.CodecType != "video" {
			continue
		}
		log.Printf("sample_aspect_ratio=%q has_b_frames=%d bit_rate=%s for %s", s.SAR, s.HasBFrames, s.BitRate, videoPath)
		var bitrate int64
		fmt.Sscan(s.BitRate, &bitrate)
		missingSAR := s.SAR == "" || s.SAR == "N/A" || s.SAR == "0:1"
		return missingSAR || s.HasBFrames > 0, bitrate, nil
	}
	return false, 0, fmt.Errorf("no video stream found in %s", videoPath)
}

// fix re-encodes the video with B-frames disabled targeting the original bitrate
// so iOS can play the file at roughly the same file size. Audio is stream-copied.
func (fw *FixIOSWorker) fix(videoPath string, bitrate int64) error {
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
			"c:a":       "copy",
			"movflags":  "+faststart",
		}).
		OverWriteOutput().
		Compile()

	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		_ = os.Remove(tmpPath)
		return err
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
