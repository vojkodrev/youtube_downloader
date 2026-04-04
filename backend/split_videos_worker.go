package main

import (
	"log"
	"path/filepath"
	"regexp"
	"time"
)

type SplitVideosWorker struct {
	cfg           *Config
	store         *VideoStore
	videoDuration *VideoDuration
	videoSplitter *VideoSplitter
}

func NewSplitVideosWorker(cfg *Config, store *VideoStore, videoDuration *VideoDuration, videoSplitter *VideoSplitter) *SplitVideosWorker {
	return &SplitVideosWorker{cfg: cfg, store: store, videoDuration: videoDuration, videoSplitter: videoSplitter}
}

func (sw *SplitVideosWorker) Start() {
	const splitDuration = 3 * 60 * 60
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
			if dur <= splitDuration {
				continue
			}
			if err := sw.videoSplitter.Split(videoPath, splitDuration); err != nil {
				log.Println("error splitting", v.Filename, ":", err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
