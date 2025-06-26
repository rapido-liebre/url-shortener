package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// --- Fake service used for handler tests ---

type fakeService struct {
	shortenCalled bool
	resolveCalled bool
}

func (f *fakeService) ShortenURL(ctx context.Context, longURL string, forceNew bool) (string, error) {
	f.shortenCalled = true
	return "http://localhost:8080/u/abc123", nil
}

func (f *fakeService) ResolveURL(ctx context.Context, shortID string) (string, error) {
	f.resolveCalled = true
	if shortID == "abc123" {
		return "https://example.org", nil
	}
	return "", nil
}

// --- Tests ---

func TestShorten_Success(t *testing.T) {
	svc := &fakeService{}
	h := NewHandler(svc)

	payload := map[string]interface{}{
		"long_url":  "https://example.org",
		"force_new": false,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/links/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Shorten(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var out map[string]string
	json.NewDecoder(resp.Body).Decode(&out)

	if out["short_url"] != "http://localhost:8080/u/abc123" {
		t.Errorf("unexpected short_url: %s", out["short_url"])
	}
	if !svc.shortenCalled {
		t.Errorf("expected ShortenURL to be called")
	}
}

func TestRedirect_Success(t *testing.T) {
	svc := &fakeService{}
	h := NewHandler(svc)

	r := chi.NewRouter()
	r.Get("/u/{id}", h.Redirect)

	req := httptest.NewRequest(http.MethodGet, "/u/abc123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected status 302, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "https://example.org" {
		t.Errorf("unexpected redirect location: %s", loc)
	}
	if !svc.resolveCalled {
		t.Errorf("expected ResolveURL to be called")
	}
}
