package fxparams

import (
	"github.com/reoden/go-echo-template/catalogs/internal/shared/data/dbcontext"
	"github.com/reoden/go-echo-template/pkg/core/messaging/producer"
	"github.com/reoden/go-echo-template/pkg/logger"
	"github.com/reoden/go-echo-template/pkg/otel/tracing"

	"go.uber.org/fx"
)

type ProductHandlerParams struct {
	fx.In

	Log               logger.Logger
	CatalogsDBContext *dbcontext.CatalogsGormDBContext
	RabbitmqProducer  producer.Producer
	Tracer            tracing.AppTracer
}
