package model

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
	DB *pgxpool.Pool
}
