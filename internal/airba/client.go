package airba

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"billing-service/internal/model"
)

type Client struct {
	BaseURL      string
	TerminalID   string
	User         string
	Password     string
	AccessToken  string
	SignatureKey string
	HTTPClient   *http.Client
}

// NewClient создаёт нового клиента AirbaPay
func NewClient(user, password, terminalID, baseURL string) *Client {
	return &Client{
		User:         user,
		Password:     password,
		TerminalID:   terminalID,
		BaseURL:      baseURL,
		HTTPClient:   &http.Client{},
		SignatureKey: baseURL, // на этапе инициализации ты передаёшь baseURL как последний аргумент, но по факту это SignatureKey
	}
}

// Authorize выполняет авторизацию и сохраняет access_token
func (c *Client) Authorize() error {
	url := fmt.Sprintf("%s/auth/sign-in", c.BaseURL)

	payload := map[string]string{
		"user":        c.User,
		"password":    c.Password,
		"terminal_id": c.TerminalID,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authorize, status: %s", resp.Status)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.AccessToken == "" {
		return errors.New("empty access token")
	}

	c.AccessToken = result.AccessToken
	return nil
}

// Send выполняет авторизованный запрос и записывает ответ в result
func (c *Client) Send(method, path string, body interface{}, result interface{}) error {
	fullURL := fmt.Sprintf("%s%s", c.BaseURL, path)

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("airba error: status %d", resp.StatusCode)
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

// CreatePayment создаёт новый платёж через AirbaPay
func (c *Client) CreatePayment(ctx context.Context, req model.CreatePaymentRequest) (*model.CreatePaymentResponse, error) {
	var response model.CreatePaymentResponse
	err := c.Send(http.MethodPost, "/api/v2/payments", req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	return &response, nil
}
