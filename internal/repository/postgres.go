package repository

import (
	"billing-service/internal/model"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateRole(ctx context.Context, role model.Role) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO roles (id, role) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET role = EXCLUDED.role`,
		role.ID, role.Role,
	)
	return err
}

func (r *Repository) CreatePayment(ctx context.Context, p model.Payment) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO payments (id, role, amount, quantity) VALUES ($1, $2, $3, $4)`,
		p.ID, p.Role, p.Amount, p.Quantity,
	)
	return err
}

func (r *Repository) GetAllPayments(ctx context.Context) ([]model.Payment, error) {
	rows, err := r.db.Query(ctx, `SELECT id, role, amount, quantity FROM payments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		err := rows.Scan(&p.ID, &p.Role, &p.Amount, &p.Quantity)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *Repository) GetPaymentsByID(ctx context.Context, id, role string) ([]model.Payment, error) {
	rows, err := r.db.Query(ctx, `SELECT id, role, amount, quantity FROM payments WHERE id=$1 AND role=$2`, id, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []model.Payment
	for rows.Next() {
		var p model.Payment
		err := rows.Scan(&p.ID, &p.Role, &p.Amount, &p.Quantity)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *Repository) SaveCardWebhook(ctx context.Context, card model.SavedCard) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO saved_cards (account_id, card_id, status)
		VALUES ($1, $2, $3)
	`, card.AccountID, card.CardID, card.Status)
	return err
}

