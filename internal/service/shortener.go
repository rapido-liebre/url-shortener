package service

import (
	"context"
	"math/rand"
	"time"

	"url_shortener/internal/model"
	"url_shortener/internal/repository"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortenerService struct {
	repo    repository.URLRepository
	baseURL string
}

func NewShortenerService(repo repository.URLRepository, baseURL string) *ShortenerService {
	rand.Seed(time.Now().UnixNano())
	return &ShortenerService{repo: repo, baseURL: baseURL}
}

func (s *ShortenerService) generateShortID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

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

func (s *ShortenerService) ResolveURL(ctx context.Context, shortID string) (string, error) {
	return s.repo.GetByShortID(ctx, shortID)
}
