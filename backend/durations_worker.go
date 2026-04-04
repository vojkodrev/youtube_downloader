package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type DurationsWorker struct {
	cfg           *Config
	store         *VideoStore
	filenames     *Filenames
	videoDuration *VideoDuration
}

func NewDurationsWorker(cfg *Config, store *VideoStore, filenames *Filenames, videoDuration *VideoDuration) *DurationsWorker {
	return &DurationsWorker{cfg: cfg, store: store, filenames: filenames, videoDuration: videoDuration}
}

func (dw *DurationsWorker) Start() {
	for {
		dw.store.Mutex.RLock()
		current := dw.store.Videos
		dw.store.Mutex.RUnlock()

		var wg sync.WaitGroup
		for _, v := range current {
			durationPath := filepath.Join(dw.cfg.StreamsDir, dw.filenames.Duration(v.Filename))
			videoPath := filepath.Join(dw.cfg.StreamsDir, v.Filename)
			// duration file already exists — skip unless the video was modified after it
			if dInfo, err := os.Stat(durationPath); err == nil {
				// video not newer than duration file, e.g. file was not replaced
				if vInfo, err := os.Stat(videoPath); err == nil && !vInfo.ModTime().After(dInfo.ModTime()) {
					continue
				}
			}
			wg.Go(func() {
				duration, err := dw.videoDuration.Get(videoPath)
				if err != nil {
					log.Println("error getting duration for", v.Filename, ":", err)
					return
				}
				if err := os.WriteFile(durationPath, []byte(strconv.FormatFloat(duration, 'f', -1, 64)), 0644); err != nil {
					log.Println("error saving duration for", v.Filename, ":", err)
				} else {
					log.Println("duration saved for", v.Filename)
				}
			})
		}
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
