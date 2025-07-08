package service

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"billing-service/internal/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo   *repository.Repository
	client *airba.Client
}

func NewPaymentService(repo *repository.Repository, client *airba.Client) *PaymentService {
	return &PaymentService{
		repo:   repo,
		client: client,
	}
}

type CreatePaymentInput struct {
	Role     string    `json:"role"`
	UserID   uuid.UUID `json:"user_id"`
	Amount   float64   `json:"amount"`
	Quantity int       `json:"quantity"`
}

// Создаёт платёж в AirbaPay и сохраняет в БД
func (s *PaymentService) CreateAndSavePayment(ctx context.Context, input CreatePaymentInput) (*model.CreatePaymentResponse, error) {
	invoiceID := uuid.New().String()

	req := model.CreatePaymentRequest{
		ID:          input.UserID.String(),
		Role:        model.Role(input.Role),
		Amount:      input.Amount,
		Quantity:    input.Quantity,
		Currency:    "KZT",
		InvoiceID:   invoiceID,
		Description: "Оплата тарифа",
	}

	resp, err := s.client.CreatePayment(ctx, req)
	if err != nil {
		return nil, err
	}

	// Сохраняем платёж в БД
	err = s.insertPayment(ctx, req, "created")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Сохраняет платёж в таблицу payments
func (s *PaymentService) insertPayment(ctx context.Context, req model.CreatePaymentRequest, status string) error {
	query := `
		INSERT INTO payments (invoice_id, role, user_id, amount, quantity, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := s.repo.DB.Exec(ctx, query,
		req.InvoiceID,
		req.Role,
		req.ID,
		req.Amount,
		req.Quantity,
		status,
	)
	return err
}

// Метод для сохранения платежа без интеграции с AirbaPay (например, для тестов)
func (s *PaymentService) CreatePayment(ctx context.Context, input CreatePaymentInput) error {
	return s.insertPayment(ctx, model.CreatePaymentRequest{
		ID:        input.UserID.String(),
		Role:      model.Role(input.Role),
		Amount:    input.Amount,
		Quantity:  input.Quantity,
		InvoiceID: uuid.New().String(),
	}, "created")
}

type PaymentRecord struct {
	ID        uuid.UUID `json:"id"`
	InvoiceID string    `json:"invoice_id"`
	Amount    float64   `json:"amount"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Получает историю платежей по user_id и роли
func (s *PaymentService) GetPayments(ctx context.Context, userID uuid.UUID, role string) ([]PaymentRecord, error) {
	query := `
		SELECT id, invoice_id, amount, quantity, status, created_at
		FROM payments
		WHERE user_id = $1 AND role = $2
		ORDER BY created_at DESC
	`
	rows, err := s.repo.DB.Query(ctx, query, userID, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PaymentRecord
	for rows.Next() {
		var p PaymentRecord
		if err := rows.Scan(&p.ID, &p.InvoiceID, &p.Amount, &p.Quantity, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}

// Обновляет статус платежа по invoice_id
func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, invoiceID string, status string) error {
	query := `UPDATE payments SET status = $1 WHERE invoice_id = $2`
	_, err := s.repo.DB.Exec(ctx, query, status, invoiceID)
	return err
}
