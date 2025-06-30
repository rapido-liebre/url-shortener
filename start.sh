#!/bin/sh

echo "[start.sh] Starting Nginx..."
# Nginx works as subprocess, not deamon
nginx -g 'daemon off;' &
NGINX_EXIT=$?

if [ $NGINX_EXIT -ne 0 ]; then
  echo "[start.sh] ERROR: Nginx failed to start (exit code $NGINX_EXIT)"
  exit $NGINX_EXIT
else
  echo "[start.sh] Nginx started successfully"
fi

# Give Nginx a while..
sleep 1

# Dynamic Select .env:
if [ -n "$RAILWAY_ENV" ]; then
  ENV_FILE="$RAILWAY_ENV"
else
  # default for local usage
  ENV_FILE=".env.docker"
fi

echo "[start.sh] Using env file: $ENV_FILE"
if [ -f "/app/$ENV_FILE" ]; then
  export $(grep -v '^#' "/app/$ENV_FILE" | xargs)
else
  echo "[start.sh] WARNING: Env file /app/$ENV_FILE not found"
fi

echo "[start.sh] Listing files in /app:"
ls -la /app

echo "[start.sh] PORT value: $PORT"

echo "[start.sh] Checking if Nginx is responding on port 80..."

curl -I http://localhost/ || echo "[start.sh] ERROR: Nginx not responding on port 80"

echo "[start.sh] Checking frontend build in /usr/share/nginx/html:"
ls -l /usr/share/nginx/html || echo "[start.sh] WARNING: Frontend files missing!"

echo "[start.sh] Starting Go backend on port 8081..."
exec /usr/bin/url-shortener -port=8081
