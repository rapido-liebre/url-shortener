# Build stage
FROM node:20 AS frontend
WORKDIR /app
COPY frontend/ .
RUN npm install && npm run build

# Stage 2 — Build Go backend
FROM golang:1.24 AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url_shortener ./cmd/api

# Stage 3 — Final image with Nginx + backend
FROM nginx:stable-alpine

# Copy frontend build
COPY --from=frontend /app/dist /usr/share/nginx/html

# Copy backend binary
COPY --from=backend /app/url_shortener /usr/bin/url-shortener

# Copy custom nginx config
COPY deploy/nginx.conf /etc/nginx/nginx.conf

# Copy execution script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Expose port 80
EXPOSE 80

# Run both: Nginx (in background) + Go app
CMD ["/start.sh"]
