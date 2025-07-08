package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewPostgres инициализирует подключение к базе данных PostgreSQL
func NewPostgres(url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	db, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	return db, nil
}
