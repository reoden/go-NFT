package app

import (
	"github.com/reoden/go-NFT/pkg/config/environment"
	"github.com/reoden/go-NFT/pkg/fxapp"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/user/internal/shared/configurations/user"

	"go.uber.org/fx"
)

type UserApplication struct {
	*user.UserServiceConfigurator
}

func NewUserApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *UserApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)
	return &UserApplication{
		UserServiceConfigurator: user.NewUserServiceConfigurator(app),
	}
}
