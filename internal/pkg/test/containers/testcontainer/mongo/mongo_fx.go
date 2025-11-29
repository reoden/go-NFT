package mongo

import (
	"context"
	"testing"

	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/mongodb"
)

var MongoContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *mongodb.MongoDbOptions, logger logger.Logger) (*mongodb.MongoDbOptions, error) {
		return NewMongoTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
