package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	"url_shortener/internal/handler"
	"url_shortener/internal/repository"
	"url_shortener/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	db := setupDatabase()
	defer db.Close()

	svc := setupService(db)
	r := setupRouter(svc)

	startServer(r)
}

// loadEnv loads environment variables from .env file
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment")
	}
}

// setupDatabase initializes and returns a PostgreSQL connection
func setupDatabase() *sql.DB {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	return db
}

// setupService creates the shortener service and handler
func setupService(db *sql.DB) *handler.Handler {
	repo := repository.NewPostgresRepository(db)
	baseURL := getEnv("BASE_URL", "http://localhost:8080")
	svc := service.NewShortenerService(repo, baseURL)
	return handler.NewHandler(svc)
}

// setupRouter sets up the chi router with middleware and routes
func setupRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
	}))

	r.Post("/links/shorten", h.Shorten)
	r.Get("/u/{id}", h.Redirect)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r
}

// startServer starts the HTTP server on the configured port
func startServer(handler http.Handler) {
	port := getEnv("PORT", "8080")
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// getEnv returns the value of an env variable or a fallback default
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
