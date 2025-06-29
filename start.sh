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

# Give Nginx a while..
sleep 1

echo "[start.sh] PORT value: $PORT"

echo "[start.sh] Checking if Nginx is responding on port 80..."

curl -I http://localhost/ || echo "[start.sh] ERROR: Nginx not responding on port 80"

echo "[start.sh] Checking frontend build in /usr/share/nginx/html:"
ls -l /usr/share/nginx/html || echo "[start.sh] WARNING: Frontend files missing!"

echo "[start.sh] Starting Go backend..."
exec /usr/bin/url-shortener
