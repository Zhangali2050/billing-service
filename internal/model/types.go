package model

import "github.com/google/uuid"

// Role — возможные роли пользователей, участвующих в оплате
type Role string

const (
	RoleStudent Role = "student"
	RoleParent  Role = "parent"
)

// RoleEntry — структура, определяющая привязку роли к конкретному UUID
type RoleEntry struct {
	ID   uuid.UUID `json:"id"`   // Идентификатор пользователя
	Role Role      `json:"role"` // Роль: "student" или "parent"
}
