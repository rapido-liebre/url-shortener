#!/bin/sh

# Run Nginx in the background
nginx &

# Run backend as main process (PID 1)
exec /usr/bin/url-shortener
