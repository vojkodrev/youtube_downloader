package main

import "time"

type Video struct {
	ID       string    `json:"id"`
	Filename string    `json:"filename"`
	Name     string    `json:"name"`
	Channel  string    `json:"channel"`
	Date     time.Time `json:"date"`
	Status   string    `json:"status"`
}
