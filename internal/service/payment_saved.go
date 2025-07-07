package service

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"fmt"
	"net/http"
)

type SavedPaymentService struct {
	client *airba.Client
}

func NewSavedPaymentService(client *airba.Client) *SavedPaymentService {
	return &SavedPaymentService{client: client}
}

func (s *SavedPaymentService) PayWithSavedCard(cardID string, payment model.SavedPaymentRequest) (model.SavedPaymentResponse, error) {
	endpoint := fmt.Sprintf("/payments/%s", cardID)

	var response model.SavedPaymentResponse
	err := s.client.Send(http.MethodPost, endpoint, payment, &response)
	return response, err
}
