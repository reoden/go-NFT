package dtos

import (
	dtoV1 "github.com/reoden/go-echo-template/catalogs/internal/products/dtos/v1"
	"github.com/reoden/go-echo-template/pkg/utils"
)

// https://echo.labstack.com/guide/response/
type GetProductsResponseDto struct {
	Products *utils.ListResult[*dtoV1.ProductDto]
}
