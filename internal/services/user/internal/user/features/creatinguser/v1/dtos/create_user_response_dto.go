package dtos

import (
	"github.com/reoden/go-NFT/pkg/core/serializer/json"

	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/response/
type CreateUserResponseDto struct {
	UserID uuid.UUID `json:"userID"`
}

func (c *CreateUserResponseDto) String() string {
	return json.PrettyPrint(c)
}
