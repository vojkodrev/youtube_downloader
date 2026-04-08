package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FixIOSWorker struct {
	cfg               *Config
	store             *VideoStore
	filenames         *Filenames
	videoIOSValidator *VideoIOSValidator
}

func NewFixIOSWorker(cfg *Config, store *VideoStore, filenames *Filenames, videoIOSValidator *VideoIOSValidator) *FixIOSWorker {
	return &FixIOSWorker{cfg: cfg, store: store, filenames: filenames, videoIOSValidator: videoIOSValidator}
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
			validation, err := fw.videoIOSValidator.Validate(videoPath)
			if err != nil {
				log.Println("error probing", v.Filename, ":", err)
				continue
			}
			if !validation.NeedsIOSFix {
				if err := os.WriteFile(fixPath, []byte("ok"), 0644); err != nil {
					log.Println("error writing ios fix marker for", v.Filename, ":", err)
				}
				continue
			}
			log.Println("fixing ios compatibility for", v.Filename)
			if err := fw.fix(videoPath, validation.Bitrate); err != nil {
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
			"vf":        "setsar=1:1",
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

	if err := os.Remove(videoPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, videoPath)
}
