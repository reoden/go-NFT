package constants

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type UserOperateTypeEnum string

const (
	FREEZE   UserOperateTypeEnum = "FREEZE"   // 冻结
	UNFREEZE UserOperateTypeEnum = "UNFREEZE" // 解冻
	LOGIN    UserOperateTypeEnum = "LOGIN"    // 登录
	REGISTER UserOperateTypeEnum = "REGISTER" // 注册
	ACTIVE   UserOperateTypeEnum = "ACTIVE"   // 激活
	AUTH     UserOperateTypeEnum = "AUTH"     // 实名认证
	MODIFY   UserOperateTypeEnum = "MODIFY"   // 修改信息
	LOGOUT   UserOperateTypeEnum = "LOGOUT"   // 登出
)

type UserStateEnum string

const (
	User_INIT   UserStateEnum = "创建成功"
	User_AUTH   UserStateEnum = "实名认证"
	User_ACTIVE UserStateEnum = "上链成功"
	User_FROZEN UserStateEnum = "冻结"
)

// Scan implements the Scanner interface for UserStateEnum
func (u *UserStateEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if bv, ok := value.([]byte); ok {
		*u = UserStateEnum(string(bv))
	} else if sv, ok := value.(string); ok {
		*u = UserStateEnum(sv)
	} else {
		return fmt.Errorf("cannot scan %T into UserStateEnum", value)
	}

	return nil
}

// Value implements the Valuer interface for UserStateEnum
func (u *UserStateEnum) Value() (driver.Value, error) {
	return string(*u), nil
}

type UserRoleEnum string

const (
	CUSTOMER UserRoleEnum = "普通用户"
	ARTIST   UserRoleEnum = "艺术家"
	ADMIN    UserRoleEnum = "管理员"
)

// Scan implements the Scanner interface for UserRoleEnum
func (u *UserRoleEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if bv, ok := value.([]byte); ok {
		*u = UserRoleEnum(string(bv))
	} else if sv, ok := value.(string); ok {
		*u = UserRoleEnum(sv)
	} else {
		return fmt.Errorf("cannot scan %T into UserStateEnum", value)
	}

	return nil
}

// Value implements the Valuer interface for UserRoleEnum
func (u *UserRoleEnum) Value() (driver.Value, error) {
	return string(*u), nil
}

const (
	DefaultNickNamePrefix    = "藏家_"
	RedisTokenBlackPrefixKey = "invalid:token:cache:"
)

// redis expire time duration
const (
	UserDataCacheExpireDuration = 2 * time.Hour
	CaptchaExpireDuration       = 5 * time.Minute
	UserTokenExpireDuration     = 24 * time.Hour
)
