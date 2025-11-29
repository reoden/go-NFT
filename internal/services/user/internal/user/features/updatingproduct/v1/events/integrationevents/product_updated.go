package integrationevents

import (
	"github.com/reoden/go-NFT/pkg/core/messaging/types"
	dto "github.com/reoden/go-NFT/user/internal/user/dtos/v1"

	uuid "github.com/satori/go.uuid"
)

type ProductUpdatedV1 struct {
	*types.Message
	*dto.ProductDto
}

func NewProductUpdatedV1(productDto *dto.ProductDto) *ProductUpdatedV1 {
	return &ProductUpdatedV1{
		Message:    types.NewMessage(uuid.NewV4().String()),
		ProductDto: productDto,
	}
}
