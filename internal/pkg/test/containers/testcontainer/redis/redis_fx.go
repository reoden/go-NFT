package redis

import (
	"context"
	"testing"

	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/redis"
)

var RedisContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *redis.RedisOptions, logger logger.Logger) (*redis.RedisOptions, error) {
		return NewRedisTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
