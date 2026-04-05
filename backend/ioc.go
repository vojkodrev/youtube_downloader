package main

import (
	"context"
	"os"
	"time"

	"go.uber.org/fx"
)

func CoreProviders() fx.Option {
	return fx.Options(
		fx.Provide(func() *Config { cfg := loadConfig(); return &cfg }),
		fx.Provide(func(cfg *Config) StreamsFS { return StreamsFS(os.DirFS(cfg.StreamsDir)) }),
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
		fx.Provide(NewFixFrameRateWorker),
		fx.Provide(NewGinServer),
	)
}

func NewIOC() *fx.App {
	return fx.New(
		CoreProviders(),

		fx.Invoke(func(
			lc fx.Lifecycle,
			pollVideosWorker *PollVideosWorker,
			thumbnailsWorker *ThumbnailsWorker,
			durationsWorker *DurationsWorker,
			splitVideosWorker *SplitVideosWorker,
			fixFrameRateWorker *FixFrameRateWorker,
			cleanupWorker *CleanupWorker) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go pollVideosWorker.Start()
					go func() {
						time.Sleep(5 * time.Second)
						go thumbnailsWorker.Start()
						go durationsWorker.Start()
						go splitVideosWorker.Start()
						go fixFrameRateWorker.Start()
						go cleanupWorker.Start()
					}()
					return nil
				},
			})
		}),

		fx.Invoke(func(lc fx.Lifecycle, ginServer *GinServer) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ginServer.registerRoutes()
					go ginServer.router.Run(":8080")
					return nil
				},
			})
		}),
	)
}
