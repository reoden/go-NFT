package infrastructure

import (
	"github.com/go-playground/validator"
	bloom "github.com/reoden/go-NFT/pkg/bloomfilter"
	"github.com/reoden/go-NFT/pkg/core"
	"github.com/reoden/go-NFT/pkg/grpc"
	"github.com/reoden/go-NFT/pkg/health"
	customEcho "github.com/reoden/go-NFT/pkg/http/customecho"
	"github.com/reoden/go-NFT/pkg/migration/goose"
	"github.com/reoden/go-NFT/pkg/otel/metrics"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/pkg/postgresgorm"
	"github.com/reoden/go-NFT/pkg/postgresmessaging"
	"github.com/reoden/go-NFT/pkg/redis"
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
	redis.Module,
	health.Module,
	tracing.Module,
	metrics.Module,
	bloom.Module,

	// Other provides
	fx.Provide(validator.New),
)
