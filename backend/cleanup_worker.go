package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func cleanupWorker(streamsDir string) {
	for {
		entries, err := os.ReadDir(streamsDir)
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
				if _, err := os.Stat(filepath.Join(streamsDir, base+ext)); err == nil {
					exists = true
					break
				}
			}
			if !exists {
				path := filepath.Join(streamsDir, name)
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
