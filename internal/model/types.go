package model
import "time"


type Role struct {
	ID   string `json:"id"`
	Role string `json:"role"` // student, parent и т.д.
}

type Payment struct {
	ID       string  `json:"id"`       // user ID
	Role     string  `json:"role"`     // student или parent
	Amount   float64 `json:"amount"`   // сумма
	Quantity int     `json:"quantity"` // кол-во месяцев
}

type AirbaAuthRequest struct {
	User       string `json:"user"`
	Password   string `json:"password"`
	TerminalID string `json:"terminal_id"`
	PaymentID  string `json:"payment_id,omitempty"`
}

type AirbaAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type AirbaPaymentRequest struct {
	InvoiceID string  `json:"invoice_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	AccountID string  `json:"account_id"`
	Desc      string  `json:"description"`
}

type AirbaPaymentResponse struct {
	ID         string `json:"id"`
	InvoiceID  string `json:"invoice_id"`
	RedirectURL string `json:"redirect_url"`
}

type AirbaAddCardRequest struct {
	AccountID string `json:"account_id"`
}

type AirbaCardResponse struct {
	ID         string `json:"id"`
	AccountID  string `json:"account_id"`
	MaskedPAN  string `json:"masked_pan"`
	Name       string `json:"name"`
	Expire     string `json:"expire"`
}

type AirbaWebhookPayload struct {
	ID          string  `json:"id"`
	InvoiceID   string  `json:"invoice_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	Signature   string  `json:"signature"`
}

const (
	StatusPaid     = "PAID"
	StatusFailed   = "FAILED"
	StatusReversed = "REVERSED"
)

type SavedCard struct {
	ID        int       `json:"id"`
	AccountID string    `json:"account_id"`
	CardID    string    `json:"card_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}


type SavedPaymentRequest struct {
    AccountID string `json:"account_id"`
    CardID    string `json:"card_id"`
    Status    string `json:"status"`
}

type SavedPaymentResponse struct {
    Success bool `json:"success"`
}

type RefundRequest struct {
    PaymentID string  `json:"payment_id"`
    Amount    float64 `json:"amount"`
    Reason    string  `json:"reason"`
}

type RefundResponse struct {
    Success   bool   `json:"success"`
    RefundID  string `json:"refund_id,omitempty"`
    ErrorText string `json:"error,omitempty"`
}

type PaymentStatusResponse struct {
    PaymentID string `json:"payment_id"`
    Status    string `json:"status"`
}