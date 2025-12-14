package models

import (
	"time"

	"github.com/reoden/go-NFT/user/internal/shared/constants"
	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	Id            int64                   `json:"id,omitempty"`
	UserId        uuid.UUID               `json:"user_id,omitempty"`
	Nickname      string                  `json:"nickname,omitempty"`
	Phone         string                  `json:"phone,omitempty"`
	State         constants.UserStateEnum `json:"state,omitempty"`
	Certification bool                    `json:"certification,omitempty"`
	RealName      string                  `json:"real_name,omitempty"`
	IdCardNo      string                  `json:"id_card_no,omitempty"`
	UserRole      constants.UserRoleEnum  `json:"user_role,omitempty"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}
