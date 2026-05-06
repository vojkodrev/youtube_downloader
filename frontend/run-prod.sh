#!/bin/bash
set -e

cd "$(dirname "$0")"

# Build
echo "Building frontend..."
npm run build

# Stop existing nginx if running
if pgrep -x nginx > /dev/null; then
    echo "Stopping existing nginx..."
    nginx -p "$(pwd)" -c nginx/nginx.conf -e nginx/nginx.log -s stop
    sleep 2
fi

# Start nginx
echo "Starting nginx on http://localhost ..."
nginx -p "$(pwd)" -c nginx/nginx.conf -e nginx/nginx.log
echo "Finished"
