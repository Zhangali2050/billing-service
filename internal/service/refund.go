package service

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"context"
	"fmt"
)

type RefundService struct {
	client *airba.Client
}

func NewRefundService(client *airba.Client) *RefundService {
	return &RefundService{client: client}
}

// Выполняет возврат средств по ID платежа и сумме
func (s *RefundService) Refund(ctx context.Context, paymentID string, amount float64) error {
	req := model.RefundRequest{
		Amount: amount,
	}

	endpoint := fmt.Sprintf("/return/%s", paymentID)
	return s.client.Send("POST", endpoint, req, nil)
}
