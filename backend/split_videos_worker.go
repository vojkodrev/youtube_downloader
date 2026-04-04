package main

import (
	"log"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

func splitVideosWorker(streamsDir string, videos *[]Video, videosMutex *sync.RWMutex) {
	const splitDuration = 3 * 60 * 60
	for {
		videosMutex.RLock()
		current := *videos
		videosMutex.RUnlock()

		splitPartRe := regexp.MustCompile(` part\d{2}\.mp4$`)
		for _, v := range current {
			if v.Status != "Ready" || splitPartRe.MatchString(v.Filename) {
				continue
			}
			videoPath := filepath.Join(streamsDir, v.Filename)
			dur, err := videoDuration(videoPath)
			if err != nil {
				log.Println("error probing", v.Filename, ":", err)
				continue
			}
			if dur <= splitDuration {
				continue
			}
			if err := splitVideo(videoPath, splitDuration); err != nil {
				log.Println("error splitting", v.Filename, ":", err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
