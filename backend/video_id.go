package main

import "github.com/google/uuid"

type VideoID struct{}

func NewVideoID() *VideoID {
	return &VideoID{}
}

func (v *VideoID) FromFilename(filename string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(filename)).String()
}
