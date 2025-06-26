package service

import (
	"context"
	"testing"
	"url_shortener/internal/model"
)

// fakeRepo is mocked repo for tests
type fakeRepo struct {
	urls map[string]string
}

func (r *fakeRepo) Save(_ context.Context, url model.URL) error {
	r.urls[url.ShortID] = url.LongURL
	return nil
}

func (r *fakeRepo) GetByShortID(_ context.Context, shortID string) (string, error) {
	url, ok := r.urls[shortID]
	if !ok {
		return "", nil
	}
	return url, nil
}

func (r *fakeRepo) FindByLongURL(_ context.Context, longURL string) (string, error) {
	for id, url := range r.urls {
		if url == longURL {
			return id, nil
		}
	}
	return "", nil
}

// --- Tests ---

func TestShortenURL_New(t *testing.T) {
	repo := &fakeRepo{urls: make(map[string]string)}
	svc := NewShortenerService(repo, "http://localhost:8080")

	short, err := svc.ShortenURL(context.Background(), "https://example.org", true)
	if err != nil {
		t.Fatal(err)
	}
	if short == "" {
		t.Error("expected non-empty short URL")
	}
}

func TestShortenURL_Existing(t *testing.T) {
	repo := &fakeRepo{urls: map[string]string{"abc123": "https://example.org"}}
	svc := NewShortenerService(repo, "http://localhost:8080")

	short, err := svc.ShortenURL(context.Background(), "https://example.org", false)
	if err != nil {
		t.Fatal(err)
	}
	if short != "http://localhost:8080/u/abc123" {
		t.Errorf("unexpected short URL: %s", short)
	}
}

func TestResolveURL(t *testing.T) {
	repo := &fakeRepo{urls: map[string]string{"abc123": "https://example.org"}}
	svc := NewShortenerService(repo, "http://localhost:8080")

	url, err := svc.ResolveURL(context.Background(), "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://example.org" {
		t.Errorf("expected https://example.org, got %s", url)
	}
}
