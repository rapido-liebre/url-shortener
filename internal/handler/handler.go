package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"url_shortener/internal/service"
)

type Handler struct {
	svc *service.ShortenerService
}

func NewHandler(svc *service.ShortenerService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LongURL  string `json:"long_url"`
		ForceNew bool   `json:"force_new"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !strings.HasPrefix(req.LongURL, "http") {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	shortURL, err := h.svc.ShortenURL(r.Context(), req.LongURL, req.ForceNew)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"short_url": shortURL})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	longURL, err := h.svc.ResolveURL(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}
