package datamodels

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/reoden/go-NFT/user/internal/shared/constants"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// https://gorm.io/docs/conventions.html
// https://gorm.io/docs/models.html#gorm-Model

// UserDataModel data model
type UserDataModel struct {
	Id            int64     `gorm:"primaryKey"`
	UserId        uuid.UUID `gorm:"column:user_id"`
	Nickname      string
	Phone         string
	State         constants.UserStateEnum
	Certification bool
	RealName      string                 `gorm:"column:real_name"`
	IdCardNo      string                 `gorm:"id_card_no"`
	UserRole      constants.UserRoleEnum `gorm:"column:user_role"`
	CreatedAt     time.Time              `gorm:"default:current_timestamp"`
	UpdatedAt     time.Time
	// for soft delete - https://gorm.io/docs/delete.html#Soft-Delete
	gorm.DeletedAt
}

// TableName overrides the table name used by UserDataModel to `user` - https://gorm.io/docs/conventions.html#TableName
func (p *UserDataModel) TableName() string {
	return "users"
}

func (p *UserDataModel) String() string {
	j, _ := json.Marshal(p)

	return string(j)
}
