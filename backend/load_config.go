package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func loadConfig() Config {
	godotenv.Load()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("could not read config.yaml:", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal("could not parse config.yaml:", err)
	}

	cfg.StreamsDir = os.Getenv("STREAMS_DIR")
	if cfg.StreamsDir == "" {
		log.Fatal("STREAMS_DIR environment variable is required")
	}

	if _, err := os.Stat(cfg.StreamsDir); os.IsNotExist(err) {
		log.Fatal("streams directory does not exist: ", cfg.StreamsDir)
	}

	return cfg
}
