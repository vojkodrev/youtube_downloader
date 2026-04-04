package main

import "go.uber.org/fx"

func NewIOC() *fx.App {
	return fx.New(
		fx.Provide(func() *Config { cfg := loadConfig(); return &cfg }),
		fx.Provide(NewVideoStore),
		fx.Provide(NewCleanupWorker),
		fx.Provide(NewVideoReader),
		fx.Provide(NewPollVideosWorker),
	)
}
