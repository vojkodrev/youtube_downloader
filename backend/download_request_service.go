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
		strings.Contains(host, "twitch.tv") ||
		strings.Contains(host, "rumble.com")
}

func (s *DownloadRequestService) ExtractPlaylistID(rawURL string) (service string, id string, err error) {
	u, parseErr := url.Parse(rawURL)
	if parseErr != nil || u.Host == "" {
		return "", "", fmt.Errorf("invalid url")
	}

	host := strings.ToLower(u.Host)
	if !strings.Contains(host, "youtube.com") {
		return "", "", fmt.Errorf("only YouTube playlists are supported")
	}

	id = u.Query().Get("list")
	if id == "" {
		return "", "", fmt.Errorf("could not determine playlist id from url")
	}

	safeIDRe := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	id = safeIDRe.ReplaceAllString(id, "_")
	return "youtube", id, nil
}

func (s *DownloadRequestService) ExtractVideoID(rawURL string) (service string, id string, err error) {
	u, parseErr := url.Parse(rawURL)
	if parseErr != nil || u.Host == "" {
		return "", "", fmt.Errorf("invalid url")
	}

	host := strings.ToLower(u.Host)
	if strings.Contains(host, "youtube.com") {
		service = "youtube"
		id = u.Query().Get("v")
	} else if strings.Contains(host, "youtu.be") {
		service = "youtube"
		id = strings.Trim(u.Path, "/")
	} else if strings.Contains(host, "twitch.tv") {
		service = "twitch"
	} else if strings.Contains(host, "rumble.com") {
		service = "rumble"
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
		return "", "", fmt.Errorf("could not determine video id from url")
	}

	safeIDRe := regexp.MustCompile(`[^a-zA-Z0-9_\-]`)
	id = safeIDRe.ReplaceAllString(id, "_")
	return service, id, nil
}
