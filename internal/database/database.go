package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database url: %w", err)

	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)

	}
	return db, nil
}
