package dtos

import dtoV1 "github.com/reoden/go-echo-template/catalogs/internal/products/dtos/v1"

// https://echo.labstack.com/guide/response/
type GetProductByIdResponseDto struct {
	Product *dtoV1.ProductDto `json:"product"`
}
