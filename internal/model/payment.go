package model
import "time"


// Запрос на создание платежа (внутренний)
type CreatePaymentRequest struct {
	ID          string  `json:"id"`           // UUID для БД
	Role        Role    `json:"role"`         // "student" или "parent"
	Amount      float64 `json:"amount"`       // сумма
	Quantity    int     `json:"quantity"`     // количество
	Currency    string  `json:"currency"`     // валюта (например, "KZT")
	InvoiceID   string  `json:"invoice_id"`   // ID от AirbaPay
	AccountID   string  `json:"account_id"`   // идентификатор пользователя у AirbaPay
	Description string  `json:"description"`  // описание
}

// Ответ от AirbaPay на создание платежа
type CreatePaymentResponse struct {
	ID          string `json:"id"`
	InvoiceID   string `json:"invoice_id"`
	RedirectURL string `json:"redirect_url"`
}

// Статус платежа от AirbaPay (используется в Webhook)
type PaymentStatusResponse struct {
	ID          string `json:"id"`
	InvoiceID   string `json:"invoice_id"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type AccessResponse struct {
	Count  int       `json:"count"`
	Amount float64   `json:"amount"`
	Until  time.Time `json:"until"`
}