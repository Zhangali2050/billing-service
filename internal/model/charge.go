package model

type ChargeRequest struct {
	Amount float64 `json:"amount"` // Сумма списания
	Role   Role    `json:"role"`   // Роль пользователя (student/parent)
	UserID string  `json:"user_id"`// UUID пользователя (строка, потом парсится)
	Token  string  `json:"token"`  // Токен карты, полученный от AirbaPay
}

type ChargeResponse struct {
	ID string `json:"id"` // ID созданного платежа
}
