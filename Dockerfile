# Stage 1: Build frontend
FROM node:20 AS frontend
WORKDIR /app
COPY frontend/ .
RUN npm install && npm run build

# Stage 2: Build backend (Go)
FROM golang:1.24 AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/api

# Stage 3: Final image with Nginx and Go app
FROM nginx:stable-alpine

# Copy frontend build to Nginx public folder
COPY --from=frontend /app/dist /usr/share/nginx/html

# Copy Go backend binary
COPY --from=backend /app/url-shortener /usr/bin/url-shortener

# Copy Nginx config
COPY deploy/nginx.conf /etc/nginx/nginx.conf

# Copy startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Port exposed to Railway
EXPOSE 80

# Listening port
ENV PORT=8080

# Run both services
ENTRYPOINT ["/start.sh"]
