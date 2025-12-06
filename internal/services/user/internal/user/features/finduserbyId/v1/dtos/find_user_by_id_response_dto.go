package dtos

import (
    "github.com/reoden/go-NFT/pkg/core/serializer/json"
    dtosv1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"
)

// https://echo.labstack.com/guide/response/
type FindUserByIdResponseDto struct {
    User *dtosv1.UserDto `json:"user"`
}

func (c *FindUserByIdResponseDto) String() string {
    return json.PrettyPrint(c)
}
