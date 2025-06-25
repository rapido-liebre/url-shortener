package repository

import (
	"context"
	"database/sql"
	"errors"
	"url_shortener/internal/model"
)

type URLRepository interface {
	Save(ctx context.Context, url model.URL) error
	Get(ctx context.Context, shortID string) (string, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, url model.URL) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO urls (short_id, long_url) VALUES ($1, $2) ON CONFLICT (short_id) DO NOTHING",
		url.ShortID, url.LongURL)
	return err
}

func (r *PostgresRepository) Get(ctx context.Context, shortID string) (string, error) {
	var longURL string
	err := r.db.QueryRowContext(ctx,
		"SELECT long_url FROM urls WHERE short_id = $1", shortID).Scan(&longURL)
	if err == sql.ErrNoRows {
		return "", errors.New("not found")
	}
	return longURL, err
}
