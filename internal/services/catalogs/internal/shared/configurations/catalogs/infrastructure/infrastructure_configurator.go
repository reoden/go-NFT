package infrastructure

import (
	"github.com/reoden/go-NFT/pkg/fxapp/contracts"
	"github.com/reoden/go-NFT/pkg/logger"
	loggingpipelines "github.com/reoden/go-NFT/pkg/logger/pipelines"
	"github.com/reoden/go-NFT/pkg/otel/metrics"
	metricspipelines "github.com/reoden/go-NFT/pkg/otel/metrics/mediatr/pipelines"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	tracingpipelines "github.com/reoden/go-NFT/pkg/otel/tracing/mediatr/pipelines"

	"github.com/mehdihadeli/go-mediatr"
)

type InfrastructureConfigurator struct {
	contracts.Application
}

func NewInfrastructureConfigurator(
	app contracts.Application,
) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: app,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(
		func(l logger.Logger, tracer tracing.AppTracer, metrics metrics.AppMetrics) error {
			err := mediatr.RegisterRequestPipelineBehaviors(
				loggingpipelines.NewMediatorLoggingPipeline(l),
				tracingpipelines.NewMediatorTracingPipeline(
					tracer,
					tracingpipelines.WithLogger(l),
				),
				metricspipelines.NewMediatorMetricsPipeline(
					metrics,
					metricspipelines.WithLogger(l),
				),
			)

			return err
		},
	)
}
