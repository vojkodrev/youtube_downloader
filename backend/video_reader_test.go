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

func TestGetVideos_FormatSegmentMp4_IsSkipped(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.f140.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	assert.Empty(t, videos)
}

func TestGetVideos_TempMp4WithFinalMp4_TempIsSkipped(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4":      &fstest.MapFile{},
		"video.temp.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "video.mp4", videos[0].Filename)
}

func TestGetVideos_Mp4WithTempMp4_ReturnsStatusProcessing(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4":      &fstest.MapFile{},
		"video.temp.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "Processing", videos[0].Status)
}

func TestGetVideos_PartXxFileWithSourceMp4_PartIsSkipped(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4":         &fstest.MapFile{},
		"video part01.mp4":  &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "video.mp4", videos[0].Filename)
}

func TestGetVideos_Mp4WithPartXxFiles_ReturnsStatusProcessing(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4":        &fstest.MapFile{},
		"video part01.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "Processing", videos[0].Status)
}

func TestGetVideos_FragmentPartFile_IsSkipped(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"[Channel] Video Title.f140.mp4.part-Frag12899.part": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	assert.Empty(t, videos)
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
