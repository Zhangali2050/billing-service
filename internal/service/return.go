package service

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"fmt"
	"net/http"
)

type ReturnService struct {
	client *airba.Client
}

func NewReturnService(client *airba.Client) *ReturnService {
	return &ReturnService{client: client}
}

func (s *ReturnService) Refund(paymentID string, req model.RefundRequest) (model.RefundResponse, error) {
	endpoint := fmt.Sprintf("/payments/%s/return", paymentID)

	var response model.RefundResponse
	err := s.client.Send(http.MethodPost, endpoint, req, &response)
	return response, err
}
