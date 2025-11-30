package integrationevents

import (
	"github.com/reoden/go-NFT/pkg/core/messaging/types"
	dtoV1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

type UserCreatedV1 struct {
	*types.Message
	*dtoV1.UserDto
}

func NewUserCreatedV1(productDto *dtoV1.UserDto) *UserCreatedV1 {
	return &UserCreatedV1{
		UserDto: productDto,
		Message: types.NewMessage(uuid.NewV4().String()),
	}
}
