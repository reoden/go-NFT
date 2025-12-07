package dtos

import (
	"github.com/reoden/go-NFT/pkg/core/serializer/json"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/response/
type LoginUserResponseDto struct {
	UserId uuid.UUID
}

func (c *LoginUserResponseDto) String() string {
	return json.PrettyPrint(c)
}
