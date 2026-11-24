package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

func Open(ctx context.Context, relativePath string) (*Database, error) {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("find home directory: %w", err)
		}
		stateHome = filepath.Join(home, ".local", "state")
	}
	path := filepath.Join(stateHome, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	u := url.URL{
		Scheme:   "file",
		Path:     path,
		RawQuery: "_pragma=busy_timeout(5000)",
	}
	db, err := sql.Open("sqlite", u.String())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(1)
	database := &Database{DB: db}

	if err := database.PingContext(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return database, nil
}

func (d *Database) Init(ctx context.Context, schema string) error {
	if _, err := d.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	return nil
}

func (d *Database) Record(ctx context.Context, query string, args ...any) error {
	_, err := d.ExecContext(ctx, query, args...)
	return err
}

func (d *Database) GetLatestSuccess(ctx context.Context, query string, args ...any) (string, error) {
	res := ""
	err := d.QueryRowContext(ctx, query, args...).Scan(&res)
	if err != nil {
		return "", err
	}
	return res, nil
}
