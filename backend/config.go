package main

type Config struct {
	StreamsDir    string  `yaml:"streams_dir"`
	SplitDuration float64 `yaml:"split_duration"`
}
