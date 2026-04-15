package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type FixIOSWorker struct {
	cfg               *Config
	store             *VideoStore
	filenames         *Filenames
	videoIOSValidator *VideoIOSValidator
	iosVideoRecoder   *IOSVideoRecoder
}

func NewFixIOSWorker(cfg *Config, store *VideoStore, filenames *Filenames, videoIOSValidator *VideoIOSValidator, iosVideoRecoder *IOSVideoRecoder) *FixIOSWorker {
	return &FixIOSWorker{cfg: cfg, store: store, filenames: filenames, videoIOSValidator: videoIOSValidator, iosVideoRecoder: iosVideoRecoder}
}

func (fw *FixIOSWorker) Start() {
	for {
		if !fw.processNext() {
			time.Sleep(1 * time.Minute)
		}
	}
}

// processNext finds the first video that needs iOS fix processing, processes it, and returns true.
// Returns false if nothing needed processing.
func (fw *FixIOSWorker) processNext() bool {
	fw.store.Mutex.RLock()
	current := fw.store.Videos
	fw.store.Mutex.RUnlock()

	for _, v := range current {
		if v.Status != "Ready" {
			continue
		}
		videoPath := filepath.Join(fw.cfg.StreamsDir, v.Filename)
		fixPath := filepath.Join(fw.cfg.StreamsDir, fw.filenames.IOSFix(v.Filename))
		// don't process if marker exists
		// process if video is newer than marker
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
			return true
		}
		log.Println("fixing ios compatibility for", v.Filename)
		if err := fw.iosVideoRecoder.Recode(videoPath, validation.Bitrate); err != nil {
			log.Println("error fixing", v.Filename, ":", err)
			return true
		}
		if err := os.WriteFile(fixPath, []byte("fixed sar"), 0644); err != nil {
			log.Println("error writing ios fix marker for", v.Filename, ":", err)
		}
		log.Println("ios fix done for", v.Filename)
		return true
	}
	return false
}
