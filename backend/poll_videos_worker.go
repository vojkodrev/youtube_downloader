package main

import (
	"log"
	"sync"
	"time"
)

func pollVideosWorker(streamsDir string, videos *[]Video, videosMap *map[string]Video, videosMutex *sync.RWMutex) {
	for {
		fetched, err := getVideos(streamsDir)
		if err != nil {
			log.Println("error fetching videos:", err)
		} else {
			m := make(map[string]Video, len(fetched))
			for _, v := range fetched {
				m[v.ID] = v
			}
			videosMutex.Lock()
			*videos = fetched
			*videosMap = m
			videosMutex.Unlock()
			log.Println("loaded", len(fetched), "videos")
		}
		time.Sleep(1 * time.Minute)
	}
}
