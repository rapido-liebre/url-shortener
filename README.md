# URL Shortener Microservice

A lightweight URL shortener built with Go, PostgreSQL and React.

## 🧩 Features

- REST API to shorten and redirect URLs
- PostgreSQL as backend storage
- React-based frontend
- Dockerized for local and cloud deployment

---

## 🚀 Running locally (with Docker)

```bash
docker-compose up --build
```

## 🔗 Local services

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432 (user: postgres, password: postgres, db: shortener)

---

## 📡 API Endpoints

### POST `/links/shorten`
Shortens a URL.

#### Request:

```json
{
  "long_url": "https://example.org"
}
```
#### Response:

```json
{
"short_url": "http://localhost:8080/u/cBzml9"
}
```

### GET /u/{id}
Redirects to the original long URL.

---

## 🗃️ Database Schema
Table created automatically on first run:

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_id VARCHAR(16) UNIQUE NOT NULL,
    long_url TEXT NOT NULL
);
```

## 🛠 Makefile Commands

```bash
make up            # Run all containers
make down          # Stop and remove containers + volume
make logs          # Tail backend logs
make curl          # Test POST /links/shorten with curl
make psql          # Access PostgreSQL CLI inside container
make frontend-dev  # Start frontend locally (npm run dev)
```

## 🌍 Deployment
Project ready to deploy this stack on Railway platform.