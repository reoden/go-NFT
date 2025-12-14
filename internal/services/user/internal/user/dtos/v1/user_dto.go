package v1

import (
	"time"

	"github.com/reoden/go-NFT/user/internal/shared/constants"
	uuid "github.com/satori/go.uuid"
)

type UserDto struct {
	Id            int64                   `json:"id"`
	UserId        uuid.UUID               `json:"user_id"`
	Nickname      string                  `json:"nickname"`
	Phone         string                  `json:"phone"`
	State         constants.UserStateEnum `json:"state"`
	Certification bool                    `json:"certification"`
	RealName      string                  `json:"real_name"`
	IdCardNo      string                  `json:"id_card_no"`
	UserRole      constants.UserRoleEnum  `json:"user_role"`
	CreatedAt     time.Time               `json:"createdAt"`
	UpdatedAt     time.Time               `json:"updatedAt"`
}
