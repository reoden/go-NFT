package v1

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type UserDto struct {
	Id        int64     `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
