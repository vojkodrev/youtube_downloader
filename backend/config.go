package main

type Config struct {
	StreamsDir      string
	SplitDuration   float64  `yaml:"split_duration"`
	RecodeQualities []string `yaml:"recode_qualities"`
}
