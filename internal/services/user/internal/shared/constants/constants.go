package constants

import "time"

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
