package app

import (
	"github.com/reoden/go-echo-template/catalogs/internal/shared/configurations/catalogs"
	"github.com/reoden/go-echo-template/pkg/config/environment"
	"github.com/reoden/go-echo-template/pkg/fxapp"
	"github.com/reoden/go-echo-template/pkg/logger"

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
