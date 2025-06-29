#!/bin/sh
nginx &
exec /usr/bin/url-shortener
