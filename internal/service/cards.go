package service

import (
	"billing-service/internal/airba"
	"fmt"
	"net/http"
)

type CardService struct {
	client *airba.Client
}

func NewCardService(client *airba.Client) *CardService {
	return &CardService{client: client}
}

func (s *CardService) AddCard(accountID string) (string, error) {
	endpoint := fmt.Sprintf("/cards/")
	body := map[string]string{
		"account_id": accountID,
	}
	var response struct {
		RedirectURL string `json:"redirect_url"`
	}
	err := s.client.Send(http.MethodPost, endpoint, body, &response)
	if err != nil {
		return "", err
	}
	return response.RedirectURL, nil
}

func (s *CardService) GetCards(accountID string) ([]map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/cards/%s", accountID)
	var response []map[string]interface{}
	err := s.client.Send(http.MethodGet, endpoint, nil, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *CardService) DeleteCard(cardID string) error {
	endpoint := fmt.Sprintf("/cards/%s", cardID)
	return s.client.Send(http.MethodDelete, endpoint, nil, nil)
}

func (s *CardService) SaveCard(ctx context.Context, card model.SavedCard) error {
	return s.repo.SaveCardWebhook(ctx, card)
}

