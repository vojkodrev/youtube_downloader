package main

import "sync"

type VideoStore struct {
	Videos    []Video
	VideosMap map[string]Video
	Mutex     sync.RWMutex
}

func NewVideoStore() *VideoStore {
	return &VideoStore{
		VideosMap: make(map[string]Video),
	}
}
