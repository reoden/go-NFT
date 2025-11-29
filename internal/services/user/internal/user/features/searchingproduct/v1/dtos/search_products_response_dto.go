package dtos

import (
	"github.com/reoden/go-NFT/pkg/utils"
	dtoV1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
