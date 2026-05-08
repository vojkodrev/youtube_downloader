package main

import (
	"log"
	"os"
	"strings"

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

	cfg.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.BackendURL = os.Getenv("BACKEND_URL")
	cfg.FrontendURL = os.Getenv("FRONTEND_URL")

	if publicHosts := os.Getenv("PUBLIC_HOSTS"); publicHosts != "" {
		cfg.PublicHosts = strings.Split(publicHosts, ",")
	}

	return cfg
}
