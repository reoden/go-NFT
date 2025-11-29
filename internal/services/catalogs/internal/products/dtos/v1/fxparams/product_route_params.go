package fxparams

import (
	"github.com/reoden/go-NFT/catalogs/internal/shared/contracts"
	"github.com/reoden/go-NFT/pkg/logger"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type ProductRouteParams struct {
	fx.In

	CatalogsMetrics *contracts.CatalogsMetrics
	Logger          logger.Logger
	ProductsGroup   *echo.Group `name:"product-echo-group"`
	Validator       *validator.Validate
}
