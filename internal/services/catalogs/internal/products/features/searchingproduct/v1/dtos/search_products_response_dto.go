package dtos

import (
	dtoV1 "github.com/reoden/go-NFT/catalogs/internal/products/dtos/v1"
	"github.com/reoden/go-NFT/pkg/utils"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
