package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Shortener interface {
	ShortenURL(ctx context.Context, longURL string, forceNew bool) (string, error)
	ResolveURL(ctx context.Context, shortID string) (string, error)
}

// Handler maps incoming HTTP requests to business logic operations
type Handler struct {
	svc Shortener // Business logic for shortening and resolving URLs
}

// NewHandler creates a new HTTP handler with the given service
func NewHandler(svc Shortener) *Handler {
	return &Handler{svc: svc}
}

// Shorten handles POST /links/shorten
// Expects a JSON payload with a long_url (and optional force_new)
// Returns a JSON object containing the generated short_url
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	log.Println("Received shorten request")

	var req struct {
		LongURL  string `json:"long_url"`
		ForceNew bool   `json:"force_new"`
	}

	// Decode and validate the request payload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !strings.HasPrefix(req.LongURL, "http") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad request"})

		return
	}

	// Generate or fetch existing short URL
	shortURL, err := h.svc.ShortenURL(r.Context(), req.LongURL, req.ForceNew)
	if err != nil {
		log.Printf("ShortenURL error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	// Return the result as JSON
	_ = json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

// Redirect handles GET /u/{id}
// Looks up the original long URL for the given short ID and redirects to it
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Resolve the long URL by short ID
	longURL, err := h.svc.ResolveURL(r.Context(), id)
	if err != nil || longURL == "" {
		http.NotFound(w, r)
		return
	}

	// Issue HTTP 302 redirect to the original URL
	http.Redirect(w, r, longURL, http.StatusFound)
}
