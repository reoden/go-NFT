package fxparams

import (
	"github.com/reoden/go-NFT/pkg/core/messaging/producer"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"

	"go.uber.org/fx"
)

type UserHandlerParams struct {
	fx.In

	Log              logger.Logger
	UserDBContext    *dbcontext.UserGormDBContext
	RabbitmqProducer producer.Producer
	Tracer           tracing.AppTracer
}
