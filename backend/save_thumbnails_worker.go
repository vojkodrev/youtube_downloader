package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func saveThumbnailsWorker(streamsDir string, videos *[]Video, videosMutex *sync.RWMutex) {
	for {
		videosMutex.RLock()
		current := *videos
		videosMutex.RUnlock()

		var wg sync.WaitGroup
		for _, v := range current {
			thumbPath := filepath.Join(streamsDir, thumbnailFilename(v.Filename))
			videoPath := filepath.Join(streamsDir, v.Filename)
			durationPath := filepath.Join(streamsDir, durationFilename(v.Filename))
			// thumbnail already exists — check if it needs to be regenerated
			if tInfo, err := os.Stat(thumbPath); err == nil {
				newerExists := false
				// regenerate if the video file was modified after the thumbnail (e.g. file was replaced)
				if vInfo, err := os.Stat(videoPath); err == nil && vInfo.ModTime().After(tInfo.ModTime()) {
					newerExists = true
				}
				// regenerate if the duration file was updated after the thumbnail (e.g. seek position changed)
				if dInfo, err := os.Stat(durationPath); err == nil && dInfo.ModTime().After(tInfo.ModTime()) {
					newerExists = true
				}
				if !newerExists {
					continue
				}
			}
			wg.Go(func() {
				if err := saveThumbnail(videoPath, thumbPath); err != nil {
					log.Println("error saving thumbnail for", v.Filename, ":", err)
				} else {
					log.Println("thumbnail generated for", v.Filename)
				}
			})
		}
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
