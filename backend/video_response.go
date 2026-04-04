package main

import "time"

type VideoResponse struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Channel string    `json:"channel"`
	Date    time.Time `json:"date"`
	Status  string    `json:"status"`
}
