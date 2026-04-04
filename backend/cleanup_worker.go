package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CleanupWorker struct {
	cfg *Config
}

func NewCleanupWorker(cfg *Config) *CleanupWorker {
	return &CleanupWorker{cfg: cfg}
}

func (cw *CleanupWorker) Start() {
	for {
		entries, err := os.ReadDir(cw.cfg.StreamsDir)
		if err != nil {
			log.Println("cleanup: error reading streams dir:", err)
			continue
		}

		for _, entry := range entries {
			name := entry.Name()
			var base string
			if b, ok := strings.CutSuffix(name, ".duration.txt"); ok {
				base = b
			} else if b, ok := strings.CutSuffix(name, ".jpg"); ok {
				base = b
			} else {
				continue
			}

			exists := false
			for _, ext := range []string{".mp4", ".part"} {
				if _, err := os.Stat(filepath.Join(cw.cfg.StreamsDir, base+ext)); err == nil {
					exists = true
					break
				}
			}
			if !exists {
				path := filepath.Join(cw.cfg.StreamsDir, name)
				if err := os.Remove(path); err != nil {
					log.Println("cleanup: error removing", name, ":", err)
				} else {
					log.Println("cleanup: removed orphaned file", name)
				}
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
