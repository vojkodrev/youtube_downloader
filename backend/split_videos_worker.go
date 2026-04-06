package main

import (
	"log"
	"path/filepath"
	"regexp"
	"time"
)

type SplitVideosWorker struct {
	cfg              *Config
	store            *VideoStore
	videoDuration    *VideoDuration
	videoSplitter    *VideoSplitter
	videoIOSValidator *VideoIOSValidator
}

func NewSplitVideosWorker(cfg *Config, store *VideoStore, videoDuration *VideoDuration, videoSplitter *VideoSplitter, videoIOSValidator *VideoIOSValidator) *SplitVideosWorker {
	return &SplitVideosWorker{cfg: cfg, store: store, videoDuration: videoDuration, videoSplitter: videoSplitter, videoIOSValidator: videoIOSValidator}
}

func (sw *SplitVideosWorker) Start() {
	if sw.cfg.SplitDuration == 0 {
		log.Fatal("split_duration is not set in config")
	}
	for {
		sw.store.Mutex.RLock()
		current := sw.store.Videos
		sw.store.Mutex.RUnlock()

		splitPartRe := regexp.MustCompile(` part\d{2}\.mp4$`)
		for _, v := range current {
			if v.Status != "Ready" || splitPartRe.MatchString(v.Filename) {
				continue
			}
			videoPath := filepath.Join(sw.cfg.StreamsDir, v.Filename)
			dur, err := sw.videoDuration.Get(videoPath)
			if err != nil {
				log.Println("error probing", v.Filename, ":", err)
				continue
			}
			if dur <= sw.cfg.SplitDuration {
				continue
			}
			validation, err := sw.videoIOSValidator.Validate(videoPath)
			if err != nil {
				log.Println("error validating ios compatibility for", v.Filename, ":", err)
				continue
			}
			if validation.NeedsIOSFix {
				log.Println("skipping split for", v.Filename, "— ios fix pending")
				continue
			}
			if err := sw.videoSplitter.Split(videoPath, sw.cfg.SplitDuration); err != nil {
				log.Println("error splitting", v.Filename, ":", err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
