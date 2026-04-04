package main

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func setupVideoReader(t *testing.T, fsys fstest.MapFS) *VideoReader {
	var vr *VideoReader
	app := fxtest.New(t,
		CoreProviders(),
		fx.Decorate(func() StreamsFS { return StreamsFS(fsys) }),
		fx.Populate(&vr),
	)
	app.RequireStart()
	t.Cleanup(func() { app.RequireStop() })
	return vr
}

func TestGetVideos_PartFile_ReturnsStatusDownloading(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video2026-04-04 17_55.mp4.part": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "Downloading", videos[0].Status)
}

func TestGetVideos_TwoPartFiles_ReturnsOnlyLargerWithStatusDownloading(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video2026-04-04 17_55.f140.webm.part": &fstest.MapFile{Data: make([]byte, 100)},
		"video2026-04-04 17_55.f251.webm.part": &fstest.MapFile{Data: make([]byte, 200)},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "video2026-04-04 17_55.f251.webm.part", videos[0].Filename)
	assert.Equal(t, "Downloading", videos[0].Status)
}

func TestGetVideos_Mp4File_ReturnsStatusReady(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "Ready", videos[0].Status)
}
