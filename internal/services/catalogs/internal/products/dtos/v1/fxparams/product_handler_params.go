package fxparams

import (
	"github.com/reoden/go-NFT/catalogs/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/pkg/core/messaging/producer"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"

	"go.uber.org/fx"
)

type ProductHandlerParams struct {
	fx.In

	Log               logger.Logger
	CatalogsDBContext *dbcontext.CatalogsGormDBContext
	RabbitmqProducer  producer.Producer
	Tracer            tracing.AppTracer
}
