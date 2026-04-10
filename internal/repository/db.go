package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

func NewDB(ctx context.Context, dbURL string) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	slog.Info("database connected successfully")

	return db, nil
}
