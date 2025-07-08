package service

import (
	"billing-service/internal/airba"
	"billing-service/internal/model"
	"context"
	"fmt"
)

type CardService struct {
	client *airba.Client
}

func NewCardService(client *airba.Client) *CardService {
	return &CardService{client: client}
}

// 1. Запросить redirect_url для добавления карты
func (s *CardService) AddCard(ctx context.Context, accountId string) (string, error) {
	req := model.SaveCardRequest{
		AccountID: accountId,
	}
	var resp model.SaveCardResponse
	err := s.client.Send("POST", "/api/v2/cards", req, &resp)
	if err != nil {
		return "", err
	}
	return resp.RedirectURL, nil
}

// 2. Получить список сохранённых карт
func (s *CardService) GetCards(ctx context.Context, accountId string) ([]model.CardInfo, error) {
	var cards []model.CardInfo
	err := s.client.Send("GET", fmt.Sprintf("/api/v2/cards/%s", accountId), nil, &cards)
	return cards, err
}

// 3. Удалить карту по ID
func (s *CardService) DeleteCard(ctx context.Context, cardID string) error {
	return s.client.Send("DELETE", fmt.Sprintf("/api/v2/cards/%s", cardID), nil, nil)
}
