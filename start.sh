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

echo "[start.sh] Checking env config..."

# Only source .env if DATABASE_URL is not already set (e.g. Railway)
if [ -z "$DATABASE_URL" ]; then
  if [ -f /app/.env ]; then
    echo "[start.sh] Sourcing .env from /app/.env"
    export $(grep -v '^#' /app/.env | xargs)
  else
    echo "[start.sh] WARNING: /app/.env not found"
  fi
else
  echo "[start.sh] DATABASE_URL already set – skipping .env loading"
fi

echo "[start.sh] DATABASE_URL=$DATABASE_URL"
echo "[start.sh] BASE_URL=$BASE_URL"
echo "[start.sh] PORT=$PORT"

echo "[start.sh] Checking if Nginx is responding on port 80..."
curl -I http://localhost/ || echo "[start.sh] ERROR: Nginx not responding on port 80"

echo "[start.sh] Checking frontend build in /usr/share/nginx/html:"
ls -l /usr/share/nginx/html || echo "[start.sh] WARNING: Frontend files missing!"

echo "[start.sh] Starting Go backend on port $PORT..."
exec /usr/bin/url-shortener -port=$PORT
