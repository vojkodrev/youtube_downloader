#!/bin/bash
set -e

cd "$(dirname "$0")"

NGINX_CONF="nginx/nginx-linux.conf"

stop_nginx() {
    if pgrep -x nginx > /dev/null; then
        echo "Stopping existing nginx..."
        if [ -f nginx/nginx.pid ]; then
            sudo nginx -p "$(pwd)" -c "$NGINX_CONF" -e nginx/nginx.log -s stop
        else
            sudo systemctl stop nginx
        fi
        sleep 2
    else
        echo "nginx is not running"
    fi
}

if [ "$1" = "stop" ]; then
    stop_nginx
    exit 0
fi

# Build
echo "Building frontend..."
npm run build

# Stop existing nginx if running
stop_nginx

# Ensure nginx temp dirs exist
mkdir -p nginx/temp/client_body nginx/temp/proxy nginx/temp/fastcgi nginx/temp/uwsgi nginx/temp/scgi

# Start nginx
echo "Starting nginx on http://localhost ..."
sudo nginx -p "$(pwd)" -c "$NGINX_CONF" -e nginx/nginx.log
echo "Finished"
