package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"url_shortener/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *PostgresRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %v", err)
	}
	repo := NewPostgresRepository(db)
	return db, mock, repo
}

func TestSave_Success(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	url := model.URL{ShortID: "abc123", LongURL: "https://example.org"}

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO urls (short_id, long_url) VALUES ($1, $2) ON CONFLICT (short_id) DO NOTHING")).
		WithArgs(url.ShortID, url.LongURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Save(context.Background(), url)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFindByLongURL_Found(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	longURL := "https://example.org"
	expectedShortID := "abc123"

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT short_id FROM urls WHERE long_url = $1")).
		WithArgs(longURL).
		WillReturnRows(sqlmock.NewRows([]string{"short_id"}).AddRow(expectedShortID))

	shortID, err := repo.FindByLongURL(context.Background(), longURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shortID != expectedShortID {
		t.Errorf("expected %s, got %s", expectedShortID, shortID)
	}
}

func TestFindByLongURL_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	longURL := "https://nonexistent.org"

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT short_id FROM urls WHERE long_url = $1")).
		WithArgs(longURL).
		WillReturnError(sql.ErrNoRows)

	shortID, err := repo.FindByLongURL(context.Background(), longURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shortID != "" {
		t.Errorf("expected empty shortID, got %s", shortID)
	}
}

func TestGetByShortID_Found(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	shortID := "abc123"
	expectedLongURL := "https://example.org"

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT long_url FROM urls WHERE short_id = $1")).
		WithArgs(shortID).
		WillReturnRows(sqlmock.NewRows([]string{"long_url"}).AddRow(expectedLongURL))

	longURL, err := repo.GetByShortID(context.Background(), shortID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if longURL != expectedLongURL {
		t.Errorf("expected %s, got %s", expectedLongURL, longURL)
	}
}

func TestGetByShortID_NotFound(t *testing.T) {
	db, mock, repo := setupMockDB(t)
	defer db.Close()

	shortID := "unknown"

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT long_url FROM urls WHERE short_id = $1")).
		WithArgs(shortID).
		WillReturnError(sql.ErrNoRows)

	longURL, err := repo.GetByShortID(context.Background(), shortID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if longURL != "" {
		t.Errorf("expected empty longURL, got %s", longURL)
	}
}
