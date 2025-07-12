package service

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"billing-service/internal/repository"
	"context"
	"time"
	"fmt"
	"strconv"
	"github.com/google/uuid"
	"database/sql"
	"errors"
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
	UserID   int64     `json:"user_id"`
	Amount   float64   `json:"amount"`
	Quantity int       `json:"quantity"`
}

// Создаёт платёж в AirbaPay и сохраняет в БД
func (s *PaymentService) CreateAndSavePayment(ctx context.Context, input CreatePaymentInput) (*model.CreatePaymentResponse, error) {
	invoiceID := uuid.New().String()

	req := model.CreatePaymentRequest{
		ID:          strconv.FormatInt(input.UserID, 10),
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
		ID:        strconv.FormatInt(input.UserID, 10),
		Role:      model.Role(input.Role),
		Amount:    input.Amount,
		Quantity:  input.Quantity,
		InvoiceID: uuid.New().String(),
	}, "created")
}

type PaymentRecord struct {
	ID        int64     `json:"id"`
	InvoiceID string    `json:"invoice_id"`
	Amount    float64   `json:"amount"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Получает историю платежей по user_id и роли
func (s *PaymentService) GetPayments(ctx context.Context, userID int64, role string) ([]PaymentRecord, error) {
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


type AccessData struct {
	UserID int64
	Role   string
	Amount float64
	Count  int
	Until  time.Time
}

// Вставка доступа в таблицу payments
func (s *PaymentService) GrantAccess(ctx context.Context, data AccessData) error {
	query := `
		INSERT INTO payments (user_id, role, invoice_id, amount, quantity, status, created_at, until)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.repo.DB.Exec(ctx, query,
		data.UserID,
		data.Role,
		fmt.Sprintf("granted_%d", time.Now().UnixNano()), // invoice_id dummy
		data.Amount,
		data.Count,
		"granted",
		time.Now(),
		data.Until,
	)

	if err != nil {
		fmt.Println("❌ INSERT error:", err)
	}


	return err
}

// Получение последнего granted-доступа
func (s *PaymentService) GetAccess(ctx context.Context, userID int64, role string) (*model.AccessResponse, error) {
	query := `
		SELECT amount, quantity, until
		FROM payments
		WHERE user_id = $1 AND role = $2 AND status = 'granted'
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := s.repo.DB.QueryRow(ctx, query, userID, role)

	var resp model.AccessResponse
	if err := row.Scan(&resp.Amount, &resp.Count, &resp.Until); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &resp, nil
}



func (s *PaymentService) CreatePaymentWithAccess(ctx context.Context, userID int64, role string, amount float64, count int, until time.Time) error {
	req := model.CreatePaymentRequest{
		ID:        strconv.FormatInt(userID, 10),
		Role:      model.Role(role),
		Amount:    amount,
		Quantity:  count,
		Currency:  "KZT",
		InvoiceID: uuid.New().String(),
		Description: fmt.Sprintf("Дарение доступа до %s", until.Format("2006-01-02")),
	}

	// сохраняем платёж в статусе granted
	err := s.insertPayment(ctx, req, "granted")
	return err
}

type AccessInfo struct {
	Count  int       `json:"count"`
	Amount float64   `json:"amount"`
	Until  time.Time `json:"until"`
}

func (s *PaymentService) GetAccessInfo(ctx context.Context, userID int64, role string) (*AccessInfo, error) {
	query := `
		SELECT quantity, amount, created_at
		FROM payments
		WHERE user_id = $1 AND role = $2 AND status = 'granted'
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := s.repo.DB.QueryRow(ctx, query, userID, role)

	var count int
	var amount float64
	var createdAt time.Time

	if err := row.Scan(&count, &amount, &createdAt); err != nil {
		return nil, err
	}

	return &AccessInfo{
		Count:  count,
		Amount: amount,
		Until:  createdAt.AddDate(0, 0, 30), // например, +30 дней от created_at
	}, nil
}
