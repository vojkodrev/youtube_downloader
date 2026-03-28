package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type Config struct {
	StreamsDir string `yaml:"streams_dir"`
}

func loadConfig() Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("could not read config.yaml:", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal("could not parse config.yaml:", err)
	}

	return cfg
}

func main() {
	cfg := loadConfig()
	log.Println("streams dir:", cfg.StreamsDir)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
