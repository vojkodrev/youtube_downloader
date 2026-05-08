#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SESSION="youtube-downloader"

# Kill existing session if running
tmux kill-session -t "$SESSION" 2>/dev/null || true

# Create new session with backend in first window
tmux new-session -d -s "$SESSION" -n backend "cd '$ROOT/backend' && bash run-prod.sh; exec bash"

# Frontend in second window
tmux new-window -t "$SESSION" -n frontend "cd '$ROOT/frontend' && bash run-prod.sh; exec bash"

# Downloader in third window
tmux new-window -t "$SESSION" -n downloader "cd '$ROOT/downloader' && bash run.sh; exec bash"

# Attach to session
tmux attach-session -t "$SESSION"
