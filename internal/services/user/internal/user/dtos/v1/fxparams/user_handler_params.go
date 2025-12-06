package fxparams

import (
	bloom "github.com/reoden/go-NFT/pkg/bloomfilter"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/user/internal/user/contracts"
)

type CreateUserHandlerParams struct {
	Log             logger.Logger
	UserDBContext   *dbcontext.UserGormDBContext
	UserRepository  contracts.UserRepository
	RedisRepository contracts.UserCacheRepository
	BloomFilter     *bloom.BloomFilterFactory
	Tracer          tracing.AppTracer
}
