package database

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/v1adhope/bar-indicators/weather/internal/entity"
	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

func Open(ctx context.Context) (*Database, error) {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("find home directory: %w", err)
		}
		stateHome = filepath.Join(home, ".local", "state")
	}
	dir := filepath.Join(stateHome, "waybar-weather-indicator")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	u := url.URL{
		Scheme:   "file",
		Path:     filepath.Join(dir, "weather_history.sqlite3"),
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

//go:embed schema.sql
var schema string

func (d *Database) Init(ctx context.Context) error {
	if _, err := d.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	return nil
}

//go:embed record.sql
var record string

func (d *Database) Record(ctx context.Context, city string, status entity.Status, res string) error {
	_, err := d.ExecContext(ctx, record,
		time.Now().UTC().Format(time.RFC3339), city, status.String(), res,
	)
	return err
}

//go:embed get_latest_success.sql
var getLatestSuccess string

func (d *Database) GetLatestSuccess(ctx context.Context, city string) (string, error) {
	res := ""
	err := d.QueryRowContext(ctx, getLatestSuccess, city, entity.StatusSuccess.String()).Scan(&res)
	if err != nil {
		return "", err
	}
	return res, nil
}
