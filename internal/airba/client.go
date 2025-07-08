package airba

import (
	"context"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

)

type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	User        string
	Password    string
	TerminalID  string
	AccessToken string
}

// NewClient инициализирует клиента AirbaPay
func NewClient(baseURL, user, password, terminalID string) *Client {
	return &Client{
		BaseURL:    baseURL,
		User:       user,
		Password:   password,
		TerminalID: terminalID,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// AuthRequest — тело запроса авторизации
type AuthRequest struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	TerminalID string `json:"terminal_id"`
}

// AuthResponse — ответ от /auth/sign-in
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// Authenticate получает access_token и сохраняет в клиент
func (c *Client) Authenticate() error {
	url := fmt.Sprintf("%s/api/v1/auth/sign-in", c.BaseURL)
	payload := AuthRequest{
		User:       c.User,
		Password:   c.Password,
		TerminalID: c.TerminalID,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("auth failed: %s", data)
	}

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	c.AccessToken = result.AccessToken
	return nil
}

// GenericPost выполняет POST-запрос с авторизацией
func (c *Client) GenericPost(path string, payload any) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("http error %d: %s", resp.StatusCode, raw)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *Client) Send(ctx context.Context, req *http.Request) (*http.Response, error) {
    return c.HTTPClient.Do(req)
}
