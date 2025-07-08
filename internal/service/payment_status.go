package service

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"context"
	"fmt"
	"net/http"
)

type PaymentStatusService struct {
	client *airba.Client
	repo   *model.Repository
}

func NewPaymentStatusService(client *airba.Client, repo *model.Repository) *PaymentStatusService {
	return &PaymentStatusService{
		client: client,
		repo:   repo,
	}
}

func (s *PaymentStatusService) GetPaymentStatus(paymentID string) (model.PaymentStatusResponse, error) {
	endpoint := fmt.Sprintf("/payments/%s/status", paymentID)

	var response model.PaymentStatusResponse
	err := s.client.Send(http.MethodGet, endpoint, nil, &response)
	return response, err
}

func (s *PaymentStatusService) UpdatePaymentStatus(ctx context.Context, invoiceID string, newStatus string) error {
	query := `UPDATE payments SET status = $1 WHERE invoice_id = $2`
	_, err := s.repo.DB.Exec(ctx, query, newStatus, invoiceID)
	return err
}
