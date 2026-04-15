package main

import "sync"

type VideoStore struct {
	Videos          []Video
	VideosMap       map[string]Video
	VideoVersionsMap map[string]VideoVersion
	Mutex           sync.RWMutex
}

func NewVideoStore() *VideoStore {
	return &VideoStore{
		VideosMap:       make(map[string]Video),
		VideoVersionsMap: make(map[string]VideoVersion),
	}
}
