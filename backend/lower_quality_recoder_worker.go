package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type LowerQualityRecoderWorker struct {
	cfg           *Config
	store         *VideoStore
	filenames     *Filenames
	recoders      map[string]LowerQualityVideoRecoder
	videoDuration *VideoDuration
}

func NewLowerQualityRecoderWorker(cfg *Config, store *VideoStore, filenames *Filenames, recoders map[string]LowerQualityVideoRecoder, videoDuration *VideoDuration) *LowerQualityRecoderWorker {
	return &LowerQualityRecoderWorker{cfg: cfg, store: store, filenames: filenames, recoders: recoders, videoDuration: videoDuration}
}

func (w *LowerQualityRecoderWorker) Start() {
	for {
		if !w.processNext() {
			time.Sleep(1 * time.Minute)
		}
	}
}

// processNext finds the first video+quality that needs recoding, processes it, and returns true.
// Returns false if nothing needed processing.
func (w *LowerQualityRecoderWorker) processNext() bool {
	w.store.Mutex.RLock()
	current := w.store.Videos
	w.store.Mutex.RUnlock()

	for _, v := range current {
		if v.Status != "Ready" {
			continue
		}
		videoPath := filepath.Join(w.cfg.StreamsDir, v.Filename)
		fixPath := filepath.Join(w.cfg.StreamsDir, w.filenames.IOSFix(v.Filename))
		// don't process if iOS fix marker is missing
		if _, err := os.Stat(fixPath); err != nil {
			continue
		}
		// skip long videos that haven't been split into parts yet — recode the parts instead
		if !splitPartRe.MatchString(v.Filename) && w.cfg.SplitDuration > 0 {
			dur, err := w.videoDuration.Get(videoPath)
			if err == nil && dur > w.cfg.SplitDuration {
				continue
			}
		}
		for _, quality := range w.cfg.RecodeQualities {
			outputPath := filepath.Join(w.cfg.StreamsDir, w.filenames.LowerQuality(v.Filename, quality))
			// run if original is newer than output, skip if up-to-date
			if outInfo, err := os.Stat(outputPath); err == nil {
				if vInfo, err := os.Stat(videoPath); err != nil || !vInfo.ModTime().After(outInfo.ModTime()) {
					continue
				}
			}
			log.Println("recoding to "+quality+":", v.Filename)
			if err := w.recoders[quality].Recode(videoPath, outputPath); err != nil {
				log.Println("error recoding to "+quality, v.Filename, ":", err)
			} else {
				log.Println(quality+" recode done for", v.Filename)
			}
			return true
		}
	}
	return false
}
