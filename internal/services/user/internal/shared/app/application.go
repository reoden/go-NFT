package app

import (
	"github.com/reoden/go-NFT/pkg/config/environment"
	"github.com/reoden/go-NFT/pkg/fxapp"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/user/internal/shared/configurations/catalogs"

	"go.uber.org/fx"
)

type CatalogsApplication struct {
	*catalogs.CatalogsServiceConfigurator
}

func NewCatalogsApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *CatalogsApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)
	return &CatalogsApplication{
		CatalogsServiceConfigurator: catalogs.NewCatalogsServiceConfigurator(app),
	}
}
