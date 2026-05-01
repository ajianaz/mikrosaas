package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx stdlib driver for database/sql
)

// Config holds database connection configuration.
type Config struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Connect creates a new database connection pool with sensible defaults.
func Connect(cfg Config) (*sqlx.DB, error) {
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 5
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 5 * time.Minute
	}

	db, err := sqlx.Connect("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// Ping verifies the database connection is alive.
func Ping(db *sqlx.DB) error {
	return db.Ping()
}
