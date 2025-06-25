package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"url_shortener/internal/handler"
	"url_shortener/internal/repository"
	"url_shortener/internal/service"
)

func main() {
	// Konfiguracja połączenia z bazą
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	// Inicjalizacja warstw
	repo := repository.NewPostgresRepository(db)
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	svc := service.NewShortenerService(repo, baseURL)
	h := handler.NewHandler(svc)

	// Routing
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/links/shorten", h.Shorten)
	r.Get("/u/{id}", h.Redirect)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
