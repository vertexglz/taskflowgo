package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"taskflow/internal/task"

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
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("ping sqlite database: %w (close failed : %v)", err, closeErr)
		}
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}
	repository := &Repository{db: db}

	if err := repository.initSchema(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize repository: %w", err)
	}
	return repository, nil

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

func (r *Repository) initSchema(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks(
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	status TEXT NOT NULL,
	priority TEXT NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL,
	version INTEGER NOT NULL DEFAULT 1);`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("initialization tasks table: %w", err)
	}
	return nil
}

func (r *Repository) Save(ctx context.Context, t task.Task) error {
	if t == nil {
		return errors.New("task cannot be nil")
	}
	query := `
	INSERT INTO tasks(
	id,
	title,
	description,
	status,
	priority,
	created_at,
	updated_at)
	VALUES (?,?,?,?,?,?,?);`

	_, err := r.db.ExecContext(
		ctx,
		query,
		string(t.ID()),
		t.Title(),
		t.Description(),
		string(t.Status()),
		string(t.Priority()),
		t.CreatedAt().Format(time.RFC3339Nano),
		t.UpdatedAt().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("save task %q: %w", t.ID(), err)
	}
	return nil
}
