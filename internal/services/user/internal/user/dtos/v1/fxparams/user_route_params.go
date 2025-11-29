package fxparams

import (
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/user/internal/shared/contracts"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type UserRouteParams struct {
	fx.In

	UserMetrics *contracts.UserMetrics
	Logger      logger.Logger
	UserGroup   *echo.Group `name:"user-echo-group"`
	Validator   *validator.Validate
}
