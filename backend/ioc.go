package main

import "go.uber.org/fx"

func NewIOC() *fx.App {
	return fx.New(
		fx.Provide(func() *Config { cfg := loadConfig(); return &cfg }),
		fx.Provide(NewFilenames),
		fx.Provide(NewVideoDuration),
		fx.Provide(NewThumbnailSaver),
		fx.Provide(NewThumbnailsWorker),
		fx.Provide(NewDurationsWorker),
		fx.Provide(NewGinSharableFileServer),
		fx.Provide(NewVideoStore),
		fx.Provide(NewCleanupWorker),
		fx.Provide(NewVideoReader),
		fx.Provide(NewPollVideosWorker),
		fx.Provide(NewVideoSplitter),
		fx.Provide(NewSplitVideosWorker),
	)
}
