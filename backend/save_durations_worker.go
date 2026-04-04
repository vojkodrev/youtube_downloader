package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func saveDurationsWorker(streamsDir string, videos *[]Video, videosMutex *sync.RWMutex) {
	for {
		videosMutex.RLock()
		current := *videos
		videosMutex.RUnlock()

		var wg sync.WaitGroup
		for _, v := range current {
			durationPath := filepath.Join(streamsDir, durationFilename(v.Filename))
			videoPath := filepath.Join(streamsDir, v.Filename)
			// duration file already exists — skip unless the video was modified after it
			if dInfo, err := os.Stat(durationPath); err == nil {
				// video not newer than duration file, e.g. file was not replaced
				if vInfo, err := os.Stat(videoPath); err == nil && !vInfo.ModTime().After(dInfo.ModTime()) {
					continue
				}
			}
			wg.Go(func() {
				duration, err := videoDuration(videoPath)
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
