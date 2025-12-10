package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// UserOperateStream model
type UserOperateStream struct {
	Id          uint64
	GMTCreate   time.Time
	GMTModified time.Time
	UserId      uuid.UUID
	Type        string
	OperateTime time.Time
	Param       string
	ExtendInfo  string
	Deleted     int
	LockVersion int
}

func (u *UserOperateStream) TableName() string {
	return "user_operate_stream"
}
