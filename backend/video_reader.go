package main

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/google/uuid"
)

type VideoReader struct {
	cfg *Config
	fs  StreamsFS
}

func NewVideoReader(cfg *Config, fsys StreamsFS) *VideoReader {
	return &VideoReader{cfg: cfg, fs: fsys}
}

func (vr *VideoReader) GetVideos() ([]Video, error) {
	entries, err := fs.ReadDir(vr.fs, ".")
	if err != nil {
		return nil, err
	}

	splitPartRe := regexp.MustCompile(`^(.+) part\d{2}\.mp4$`)
	formatRe := regexp.MustCompile(`f\d{3}\.mp4$`)
	downloadingPartRe := regexp.MustCompile(`\.f\d{3}\.[^.]+\.part$`)
	channelRe := regexp.MustCompile(`^\[([^\]]+)\] ?`)
	formatSegmentRe := regexp.MustCompile(`\.f\d{3}$`)

	// pre-scan: for each base name, find the largest part file
	largestDownloadingPart := map[string]string{} // base -> filename of largest .part file
	largestDownloadingPartSize := map[string]int64{}
	for _, entry := range entries {
		if strings.ToLower(filepath.Ext(entry.Name())) != ".part" {
			continue
		}
		base := downloadingPartRe.ReplaceAllString(entry.Name(), "")
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.Size() > largestDownloadingPartSize[base] {
			largestDownloadingPartSize[base] = info.Size()
			largestDownloadingPart[base] = entry.Name()
		}
	}

	var videos []Video
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// skip non-mp4/part files, e.g. "video.jpg", "video.mp4.duration.txt"
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".mp4" && ext != ".part" {
			continue
		}
		status := "Ready"
		if ext == ".part" {
			// only include the largest part file per base (skip smaller format segments)
			base := downloadingPartRe.ReplaceAllString(entry.Name(), "")
			if largestDownloadingPart[base] != entry.Name() {
				continue
			}
			status = "Downloading"
		}
		// skip intermediate format segments, e.g. "video.f140.mp4"
		if formatRe.MatchString(entry.Name()) {
			continue
		}
		// in-progress yt-dlp download
		if base, ok := strings.CutSuffix(entry.Name(), ".temp.mp4"); ok {
			// skip temp file if the final mp4 already exists
			if _, err := fs.Stat(vr.fs, base+".mp4"); err == nil {
				continue
			}
			status = "Processing"
		} else if ext == ".mp4" {
			// if a temp file exists alongside the final mp4, mark as Processing
			base, _ := strings.CutSuffix(entry.Name(), ".mp4")
			if _, err := fs.Stat(vr.fs, base+".temp.mp4"); err == nil {
				status = "Processing"
			}
		}
		if m := splitPartRe.FindStringSubmatch(entry.Name()); m != nil {
			// this is a partXX file — skip it if the source file still exists (splitting in progress)
			if _, err := fs.Stat(vr.fs, m[1]+".mp4"); err == nil {
				continue
			}
		} else if ext == ".mp4" {
			// this is a plain mp4 — mark as Processing if any partXX files exist (splitting in progress)
			base := strings.TrimSuffix(entry.Name(), ".mp4")
			partFileRe := regexp.MustCompile(`^` + regexp.QuoteMeta(base) + ` part\d{2}\.mp4$`)
			for _, e := range entries {
				if partFileRe.MatchString(e.Name()) {
					status = "Processing"
					break
				}
			}
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		name = strings.TrimSuffix(name, ".mp4")
		name = formatSegmentRe.ReplaceAllString(name, "")
		var channel string
		if m := channelRe.FindStringSubmatch(name); m != nil {
			channel = m[1]
			name = name[len(m[0]):]
		}
		v := Video{
			ID:       uuid.NewSHA1(uuid.NameSpaceURL, []byte(entry.Name())).String(),
			Filename: entry.Name(),
			Name:     name,
			Channel:  channel,
			Date:     info.ModTime(),
			Status:   status,
		}
		// log.Printf("video: id=%s name=%s date=%s", v.ID, v.Name, v.Date.Format("2006-01-02 15:04:05"))
		videos = append(videos, v)
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].Date.After(videos[j].Date)
	})

	return videos, nil
}
