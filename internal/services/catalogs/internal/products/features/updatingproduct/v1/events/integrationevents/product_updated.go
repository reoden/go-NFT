package integrationevents

import (
	dto "github.com/reoden/go-echo-template/catalogs/internal/products/dtos/v1"
	"github.com/reoden/go-echo-template/pkg/core/messaging/types"

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
