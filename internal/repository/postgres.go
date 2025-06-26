package repository

import (
	"context"
	"database/sql"
	"url_shortener/internal/model"
)

// URLRepository defines an interface for URL persistence operations
type URLRepository interface {
	// Save stores a short-to-long URL mapping
	Save(ctx context.Context, url model.URL) error

	// FindByLongURL returns an existing short ID for the given long URL, if it exists
	FindByLongURL(ctx context.Context, longURL string) (string, error)

	// GetByShortID returns the original long URL for the given short ID.
	GetByShortID(ctx context.Context, shortID string) (string, error)
}

// PostgresRepository implements URLRepository using PostgreSQL as the storage backend
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new instance of PostgresRepository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Save inserts a new short-to-long URL mapping into the database
// If the short_id already exists, it does nothing (conflict ignored)
func (r *PostgresRepository) Save(ctx context.Context, url model.URL) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO urls (short_id, long_url) VALUES ($1, $2) ON CONFLICT (short_id) DO NOTHING",
		url.ShortID, url.LongURL)
	return err
}

// FindByLongURL looks up the short ID for a given long URL
// Returns empty string if no match is found
func (r *PostgresRepository) FindByLongURL(ctx context.Context, longURL string) (string, error) {
	var shortID string
	err := r.db.QueryRowContext(ctx,
		"SELECT short_id FROM urls WHERE long_url = $1", longURL).Scan(&shortID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return shortID, err
}

// GetByShortID retrieves the original long URL for the given short ID
// Returns empty string if no match is found
func (r *PostgresRepository) GetByShortID(ctx context.Context, shortID string) (string, error) {
	var longURL string
	err := r.db.QueryRowContext(ctx,
		"SELECT long_url FROM urls WHERE short_id = $1", shortID).Scan(&longURL)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return longURL, err
}
