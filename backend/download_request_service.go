package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type DownloadRequestService struct{}

func NewDownloadRequestService() *DownloadRequestService {
	return &DownloadRequestService{}
}

func (s *DownloadRequestService) IsSupportedURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return false
	}
	host := strings.ToLower(u.Host)
	return strings.Contains(host, "youtube.com") ||
		strings.Contains(host, "youtu.be") ||
		strings.Contains(host, "twitch.tv")
}

func (s *DownloadRequestService) ExtractVideoID(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return "", fmt.Errorf("invalid url")
	}

	var id string

	host := strings.ToLower(u.Host)
	if strings.Contains(host, "youtube.com") {
		id = u.Query().Get("v")
	} else if strings.Contains(host, "youtu.be") {
		id = strings.Trim(u.Path, "/")
	}

	// Generic fallback: last non-empty path segment
	if id == "" {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] != "" {
				id = parts[i]
				break
			}
		}
	}

	if id == "" {
		return "", fmt.Errorf("could not determine video id from url")
	}

	safeIDRe := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	id = safeIDRe.ReplaceAllString(id, "_")
	return id, nil
}
