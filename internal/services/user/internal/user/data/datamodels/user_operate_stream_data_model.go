package datamodels

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// UserOperateStreamDataModel data model
type UserOperateStreamDataModel struct {
	Id          uint64     `gorm:"column:id;primary_key" json:"id"`
	GMTCreate   *time.Time `gorm:"column:gmt_create" json:"gmt_create"`
	GMTModified *time.Time `gorm:"column:gmt_modified" json:"gmt_modified"`
	UserId      uuid.UUID  `gorm:"column:user_id;type:varchar(64)" json:"user_id"`
	Type        string     `gorm:"column:type;type:varchar(64)" json:"type"`
	OperateTime *time.Time `gorm:"column:operate_time" json:"operate_time"`
	Param       string     `gorm:"column:param;type:text" json:"param"`
	ExtendInfo  string     `gorm:"column:extend_info;type:text" json:"extend_info"`
	Deleted     *int       `gorm:"column:deleted" json:"deleted"`
	LockVersion *int       `gorm:"column:lock_version" json:"lock_version"`
}

func (u *UserOperateStreamDataModel) TableName() string {
	return "user_operate_stream"
}
