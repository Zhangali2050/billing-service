package model

// Запрос на возврат средств
type RefundRequest struct {
	PaymentID string  `json:"payment_id"` // ID платежа, с которого возвращаются средства
	Amount    float64 `json:"amount"`     // Сумма возврата
}
