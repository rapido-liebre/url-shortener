package service

import (
	"context"
	"math/rand"
	"time"

	"url_shortener/internal/model"
	"url_shortener/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ShortenerService handles business logic for shortening and resolving URLs
type ShortenerService struct {
	repo    repository.URLRepository // Interface to the persistence layer
	baseURL string                   // Base URL used to build the full short link
}

// NewShortenerService creates a new instance of ShortenerService.
// Seeds the random generator for short ID generation
func NewShortenerService(repo repository.URLRepository, baseURL string) *ShortenerService {
	rand.Seed(time.Now().UnixNano())
	return &ShortenerService{repo: repo, baseURL: baseURL}
}

// generateShortID returns a random alphanumeric string of given length
// Used to create unique short URL identifiers
func (s *ShortenerService) generateShortID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// ShortenURL creates or retrieves a short URL for a given long URL
// If forceNew is false and the URL already exists, the existing short ID is reused
func (s *ShortenerService) ShortenURL(ctx context.Context, longURL string, forceNew bool) (string, error) {
	if !forceNew {
		existing, err := s.repo.FindByLongURL(ctx, longURL)
		if err != nil {
			return "", err
		}
		if existing != "" {
			return s.baseURL + "/u/" + existing, nil
		}
	}

	shortID := s.generateShortID(6)
	err := s.repo.Save(ctx, model.URL{ShortID: shortID, LongURL: longURL})
	if err != nil {
		return "", err
	}
	return s.baseURL + "/u/" + shortID, nil
}

// ResolveURL returns the original long URL for a given short ID
func (s *ShortenerService) ResolveURL(ctx context.Context, shortID string) (string, error) {
	return s.repo.GetByShortID(ctx, shortID)
}
