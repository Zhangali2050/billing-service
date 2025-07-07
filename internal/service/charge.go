package service

import (
	"billing-service/internal/airba"
	"fmt"
	"net/http"
)

type ChargeService struct {
	client *airba.Client
}

func NewChargeService(client *airba.Client) *ChargeService {
	return &ChargeService{client: client}
}

func (s *ChargeService) Charge(paymentID string, amount int) error {
	endpoint := fmt.Sprintf("/payments/%s/charge", paymentID)
	body := map[string]interface{}{
		"amount": amount,
	}
	return s.client.Send(http.MethodPost, endpoint, body, nil)
}
