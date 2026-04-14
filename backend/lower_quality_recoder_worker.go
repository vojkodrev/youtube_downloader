package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type LowerQualityRecoderWorker struct {
	cfg               *Config
	store             *VideoStore
	filenames         *Filenames
	videoIOSValidator *VideoIOSValidator
	recoder           *LowerQualityVideoRecoder
}

func NewLowerQualityRecoderWorker(cfg *Config, store *VideoStore, filenames *Filenames, videoIOSValidator *VideoIOSValidator, recoder *LowerQualityVideoRecoder) *LowerQualityRecoderWorker {
	return &LowerQualityRecoderWorker{cfg: cfg, store: store, filenames: filenames, videoIOSValidator: videoIOSValidator, recoder: recoder}
}

func (w *LowerQualityRecoderWorker) Start() {
	for {
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
			outputPath := filepath.Join(w.cfg.StreamsDir, w.filenames.LowerQuality720p(v.Filename))
			// don't recode if 720p already exists
			// recode if original is newer than 720p
			if outInfo, err := os.Stat(outputPath); err == nil {
				if vInfo, err := os.Stat(videoPath); err != nil || !vInfo.ModTime().After(outInfo.ModTime()) {
					continue
				}
			}
			validation, err := w.videoIOSValidator.Validate(videoPath)
			if err != nil {
				log.Println("error probing", v.Filename, ":", err)
				continue
			}
			log.Println("recoding to 720p:", v.Filename)
			if err := w.recoder.Recode(videoPath, outputPath, validation.Bitrate); err != nil {
				log.Println("error recoding to 720p", v.Filename, ":", err)
				continue
			}
			log.Println("720p recode done for", v.Filename)
		}
		time.Sleep(1 * time.Minute)
	}
}
