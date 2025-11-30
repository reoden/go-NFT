package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	Id        int64
	UserId    uuid.UUID
	Nickname  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
