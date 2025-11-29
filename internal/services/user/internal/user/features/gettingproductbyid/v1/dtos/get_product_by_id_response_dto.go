package dtos

import dtoV1 "github.com/reoden/go-NFT/user/internal/user/dtos/v1"

// https://echo.labstack.com/guide/response/
type GetProductByIdResponseDto struct {
	Product *dtoV1.ProductDto `json:"product"`
}
