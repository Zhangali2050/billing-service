package repository

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

// Repository содержит обёртку над пулом подключений к БД
type Repository struct {
	DB *pgxpool.Pool
}

// NewRepository создаёт новый репозиторий
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}
