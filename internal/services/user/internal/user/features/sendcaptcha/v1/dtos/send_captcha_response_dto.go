package dtos

import (
	"github.com/reoden/go-NFT/pkg/core/serializer/json"
)

// https://echo.labstack.com/guide/response/
type SendCaptchaResponseDto struct {
}

func (c *SendCaptchaResponseDto) String() string {
	return json.PrettyPrint(c)
}
