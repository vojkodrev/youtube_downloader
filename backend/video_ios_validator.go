package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoIOSValidation struct {
	NeedsIOSFix bool
	Bitrate     int64
}

type VideoIOSValidator struct{}

func NewVideoIOSValidator() *VideoIOSValidator {
	return &VideoIOSValidator{}
}

// Validate probes the video stream and returns iOS compatibility info and bitrate.
// NeedsIOSFix is true when the stream has B-frames or a missing SAR, either of
// which causes iOS to refuse playback.
// A working file has has_b_frames=0 and a valid sample_aspect_ratio like "1:1".
func (v *VideoIOSValidator) Validate(videoPath string) (VideoIOSValidation, error) {
	probeJSON, err := ffmpeg.ProbeWithTimeout(videoPath, 15*time.Second, ffmpeg.KwArgs{})
	if err != nil {
		return VideoIOSValidation{}, err
	}
	var probe struct {
		Streams []struct {
			CodecType  string `json:"codec_type"`
			SAR        string `json:"sample_aspect_ratio"`
			HasBFrames int    `json:"has_b_frames"`
			BitRate    string `json:"bit_rate"`
		} `json:"streams"`
	}
	if err := json.Unmarshal([]byte(probeJSON), &probe); err != nil {
		return VideoIOSValidation{}, err
	}
	for _, s := range probe.Streams {
		if s.CodecType != "video" {
			continue
		}
		log.Printf("sample_aspect_ratio=%q has_b_frames=%d bit_rate=%s for %s", s.SAR, s.HasBFrames, s.BitRate, videoPath)
		var bitrate int64
		fmt.Sscan(s.BitRate, &bitrate)
		missingSAR := s.SAR == "" || s.SAR == "N/A" || s.SAR == "0:1"
		return VideoIOSValidation{
			NeedsIOSFix: missingSAR || s.HasBFrames > 0,
			Bitrate:     bitrate,
		}, nil
	}
	return VideoIOSValidation{}, fmt.Errorf("no video stream found in %s", videoPath)
}
