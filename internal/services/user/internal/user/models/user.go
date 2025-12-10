package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	Id        int64     `json:"id,omitempty"`
	UserId    uuid.UUID `json:"user_id,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
