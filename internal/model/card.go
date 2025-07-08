package model

type SaveCardRequest struct {
	UserID    string `json:"user_id"`    // UUID пользователя (string, но потом парсится)
	Role      Role   `json:"role"`       // "student" или "parent"
	AccountID string `json:"account_id"` // ID в системе AirbaPay
}

type SaveCardResponse struct {
	RedirectURL string `json:"redirect_url"` // URL для перенаправления пользователя на форму сохранения карты
}

type CardInfo struct {
	ID       string `json:"id"`        // Уникальный ID карты (может быть нужен для удаления)
	CardMask string `json:"card_mask"` // Маскированный номер карты
	Token    string `json:"token"`     // Токен для оплаты (используется при списании)
}
