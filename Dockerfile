# Stage 1: build
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o url_shortener ./cmd/api

# Stage 2: slim runtime
FROM golang:1.24 AS runtime
WORKDIR /app
COPY --from=builder /app/url_shortener .
EXPOSE 8080
CMD ["./url_shortener"]
