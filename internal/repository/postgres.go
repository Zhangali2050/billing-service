package repository

import (
	"context"
	"database/sql"
	"billing-service/internal/model"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateRole(ctx context.Context, role model.Role) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO roles (id, role) VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, role.ID, role.Role)
	return err
}

func (r *Repository) CreatePayment(ctx context.Context, p model.Payment) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO payments (user_id, role, amount, quantity)
		VALUES ($1, $2, $3, $4)
	`, p.UserID, p.Role, p.Amount, p.Quantity)
	return err
}

func (r *Repository) GetAllPayments(ctx context.Context) ([]model.Payment, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT user_id, role, date, amount, quantity FROM payments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.UserID, &p.Role, &p.Date, &p.Amount, &p.Quantity); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}

func (r *Repository) GetPaymentsByID(ctx context.Context, id string, role string) ([]model.Payment, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT user_id, role, date, amount, quantity
		FROM payments WHERE user_id=$1 AND role=$2
	`, id, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.Payment
	for rows.Next() {
		var p model.Payment
		if err := rows.Scan(&p.UserID, &p.Role, &p.Date, &p.Amount, &p.Quantity); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}
