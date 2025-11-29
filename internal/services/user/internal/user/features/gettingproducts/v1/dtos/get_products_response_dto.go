package dtos

import (
	"github.com/reoden/go-NFT/pkg/utils"
	dtoV1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"
)

// https://echo.labstack.com/guide/response/
type GetProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
