package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type SplitVideosWorker struct {
	cfg           *Config
	store         *VideoStore
	filenames     *Filenames
	videoDuration *VideoDuration
	videoSplitter *VideoSplitter
}

func NewSplitVideosWorker(cfg *Config, store *VideoStore, filenames *Filenames, videoDuration *VideoDuration, videoSplitter *VideoSplitter) *SplitVideosWorker {
	return &SplitVideosWorker{cfg: cfg, store: store, filenames: filenames, videoDuration: videoDuration, videoSplitter: videoSplitter}
}

var splitPartRe = regexp.MustCompile(` part\d{2}\.mp4$`)

func (sw *SplitVideosWorker) Start() {
	if sw.cfg.SplitDuration == 0 {
		log.Fatal("split_duration is not set in config")
	}
	for {
		if !sw.processNext() {
			time.Sleep(1 * time.Minute)
		}
	}
}

// processNext finds the first video that needs splitting, processes it, and returns true.
// Returns false if nothing needed processing.
func (sw *SplitVideosWorker) processNext() bool {
	sw.store.Mutex.RLock()
	current := sw.store.Videos
	sw.store.Mutex.RUnlock()

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
		fixPath := filepath.Join(sw.cfg.StreamsDir, sw.filenames.IOSFix(v.Filename))
		if _, err := os.Stat(fixPath); err != nil {
			log.Println("skipping split for", v.Filename, "— ios fix not done yet")
			continue
		}
		if err := sw.videoSplitter.Split(videoPath, sw.cfg.SplitDuration); err != nil {
			log.Println("error splitting", v.Filename, ":", err)
		}
		return true
	}
	return false
}
