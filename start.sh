#!/bin/sh

echo "[start.sh] Starting Nginx..."
nginx
NGINX_EXIT=$?

if [ $NGINX_EXIT -ne 0 ]; then
  echo "[start.sh] ERROR: Nginx failed to start (exit code $NGINX_EXIT)"
  exit $NGINX_EXIT
else
  echo "[start.sh] Nginx started successfully"
fi

echo "[start.sh] Starting Go backend..."
exec /usr/bin/url-shortener
