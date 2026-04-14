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
	assert.Equal(t, "video.mp4", videos[0].Filename)
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

func TestGetVideos_RumbleDownloadingFiles_ReturnsOnlyMainPartAsDownloading(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4.part":              &fstest.MapFile{Data: make([]byte, 500)},
		"video.mp4.part-Frag154.part": &fstest.MapFile{Data: make([]byte, 100)},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "video.mp4.part", videos[0].Filename)
	assert.Equal(t, "Downloading", videos[0].Status)
}

func TestGetVideos_MixedFiles_ReturnsCorrectVideosAndStatuses(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		// ready
		"ready.mp4": &fstest.MapFile{},
		// downloading — two format parts, only larger returned
		"downloading.f140.webm.part": &fstest.MapFile{Data: make([]byte, 100)},
		"downloading.f251.webm.part": &fstest.MapFile{Data: make([]byte, 200)},
		// processing — mp4 with temp alongside
		"processing.mp4":      &fstest.MapFile{},
		"processing.temp.mp4": &fstest.MapFile{},
		// processing — mp4 with partXX files
		"splitting.mp4":        &fstest.MapFile{},
		"splitting part01.mp4": &fstest.MapFile{},
		// skipped — format segment
		"video.f140.mp4": &fstest.MapFile{},
		// skipped — fragment part file
		"[Channel] Video Title.f140.mp4.part-Frag12899.part": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 4)

	byFilename := make(map[string]Video)
	for _, v := range videos {
		byFilename[v.Filename] = v
	}

	assert.Equal(t, "Ready", byFilename["ready.mp4"].Status)
	assert.Equal(t, "Downloading", byFilename["downloading.f251.webm.part"].Status)
	assert.Equal(t, "Processing", byFilename["processing.mp4"].Status)
	assert.Equal(t, "Processing", byFilename["splitting.mp4"].Status)
}

func TestGetVideos_TwoFormatSegmentsWithTempMp4_ReturnsSingleProcessingVideo(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.f140.mp4": &fstest.MapFile{},
		"video.f251.mp4": &fstest.MapFile{},
		"video.temp.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "video.temp.mp4", videos[0].Filename)
	assert.Equal(t, "Processing", videos[0].Status)
}

func TestGetVideos_Mp4With720pAnd480p_ReturnsOriginalWithBothInVersions(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4":      &fstest.MapFile{},
		"video.720p.mp4": &fstest.MapFile{},
		"video.480p.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "video.mp4", videos[0].Filename)
	assert.ElementsMatch(t, []string{"720p", "480p"}, videos[0].Versions)
}

func TestGetVideos_Mp4With720pTempFile_ReturnsMp4AsProcessingWithNoVersions(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4":           &fstest.MapFile{},
		"video.720p.temp.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "video.mp4", videos[0].Filename)
	assert.Equal(t, "Processing", videos[0].Status)
	assert.Empty(t, videos[0].Versions)
}

func TestGetVideos_Mp4Without720p_ReturnsEmptyVersions(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Empty(t, videos[0].Versions)
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
