package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(ctx context.Context, dbPath string) (*Repository, error) {
	if strings.TrimSpace(dbPath) == "" {
		return nil, errors.New("sqlite database path cannot be empty")
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database:%w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("ping sqlite database: %w (close failed : %v)", err, closeErr)
		}
		return nil, fmt.Errorf("ping sqlite database:%w", err)
	}
	return &Repository{
		db: db,
	}, nil

}

func (r *Repository) Close() error {
	if r == nil || r.db == nil {
		return nil
	}
	if err := r.db.Close(); err != nil {
		return fmt.Errorf("close sqlite database: %w", err)
	}
	return nil
}
