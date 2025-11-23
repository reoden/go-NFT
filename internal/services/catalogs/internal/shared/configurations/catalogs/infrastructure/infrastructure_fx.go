package infrastructure

import (
	rabbitmq2 "github.com/reoden/go-echo-template/catalogs/internal/products/configurations/rabbitmq"
	"github.com/reoden/go-echo-template/pkg/core"
	"github.com/reoden/go-echo-template/pkg/grpc"
	"github.com/reoden/go-echo-template/pkg/health"
	customEcho "github.com/reoden/go-echo-template/pkg/http/customecho"
	"github.com/reoden/go-echo-template/pkg/migration/goose"
	"github.com/reoden/go-echo-template/pkg/otel/metrics"
	"github.com/reoden/go-echo-template/pkg/otel/tracing"
	"github.com/reoden/go-echo-template/pkg/postgresgorm"
	"github.com/reoden/go-echo-template/pkg/postgresmessaging"
	"github.com/reoden/go-echo-template/pkg/rabbitmq"
	"github.com/reoden/go-echo-template/pkg/rabbitmq/configurations"

	"github.com/go-playground/validator"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/
var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	postgresgorm.Module,
	postgresmessaging.Module,
	goose.Module,
	rabbitmq.ModuleFunc(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder)
			}
		},
	),
	health.Module,
	tracing.Module,
	metrics.Module,

	// Other provides
	fx.Provide(validator.New),
)
