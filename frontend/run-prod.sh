#!/bin/bash
set -e

cd "$(dirname "$0")"

# Build
echo "Building frontend..."
npm run build

# Use Linux-specific nginx config
NGINX_CONF="nginx/nginx-linux.conf"

# Stop existing nginx if running
if pgrep -x nginx > /dev/null; then
    echo "Stopping existing nginx..."
    sudo nginx -p "$(pwd)" -c "$NGINX_CONF" -e nginx/nginx.log -s stop
    sleep 2
fi

# Start nginx
echo "Starting nginx on http://localhost ..."
sudo nginx -p "$(pwd)" -c "$NGINX_CONF" -e nginx/nginx.log
echo "Finished"
