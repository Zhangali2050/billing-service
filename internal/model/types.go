package model

import "time"

type Role struct {
	ID   string `json:"id"`
	Role string `json:"role"` // "student" or "parent"
}

type Payment struct {
	ID       int       `json:"-"`
	UserID   string    `json:"id"`
	Role     string    `json:"role"`
	Date     time.Time `json:"date"`
	Amount   float64   `json:"amount"`
	Quantity int       `json:"quantity"`
}
