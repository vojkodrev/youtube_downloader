package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type LowerQualityRecoderWorker struct {
	cfg      *Config
	store    *VideoStore
	filenames *Filenames
	recoders map[string]LowerQualityVideoRecoder
}

func NewLowerQualityRecoderWorker(cfg *Config, store *VideoStore, filenames *Filenames, recoders map[string]LowerQualityVideoRecoder) *LowerQualityRecoderWorker {
	return &LowerQualityRecoderWorker{cfg: cfg, store: store, filenames: filenames, recoders: recoders}
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
		qualities := []struct {
			key        string
			outputPath string
		}{
			{"720p", filepath.Join(w.cfg.StreamsDir, w.filenames.LowerQuality720p(v.Filename))},
			{"480p", filepath.Join(w.cfg.StreamsDir, w.filenames.LowerQuality480p(v.Filename))},
		}
		for _, q := range qualities {
			// run if original is newer than output, skip if up-to-date
			if outInfo, err := os.Stat(q.outputPath); err == nil {
				if vInfo, err := os.Stat(videoPath); err != nil || !vInfo.ModTime().After(outInfo.ModTime()) {
					continue
				}
			}
			log.Println("recoding to "+q.key+":", v.Filename)
			if err := w.recoders[q.key].Recode(videoPath, q.outputPath); err != nil {
				log.Println("error recoding to "+q.key, v.Filename, ":", err)
			} else {
				log.Println(q.key+" recode done for", v.Filename)
			}
			return true
		}
	}
	return false
}
