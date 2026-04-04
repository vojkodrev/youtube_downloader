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

func TestGetVideos_StatusReady(t *testing.T) {
	vr := setupVideoReader(t, fstest.MapFS{
		"video.mp4": &fstest.MapFile{},
	})

	videos, err := vr.GetVideos()

	require.NoError(t, err)
	require.Len(t, videos, 1)
	assert.Equal(t, "Ready", videos[0].Status)
}
