package main

import (
	"log"
	"time"
)

type PollVideosWorker struct {
	videoReader *VideoReader
	store       *VideoStore
}

func NewPollVideosWorker(videoReader *VideoReader, store *VideoStore) *PollVideosWorker {
	return &PollVideosWorker{videoReader: videoReader, store: store}
}

func (pw *PollVideosWorker) Start() {
	for {
		fetched, err := pw.videoReader.GetVideos()
		if err != nil {
			log.Println("error fetching videos:", err)
		} else {
			m := make(map[string]Video, len(fetched))
			for _, v := range fetched {
				m[v.ID] = v
			}
			pw.store.Mutex.Lock()
			pw.store.Videos = fetched
			pw.store.VideosMap = m
			pw.store.Mutex.Unlock()
			log.Println("loaded", len(fetched), "videos")
		}
		time.Sleep(1 * time.Minute)
	}
}
