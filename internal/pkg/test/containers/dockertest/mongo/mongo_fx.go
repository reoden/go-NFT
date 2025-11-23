package mongo

import (
	"context"
	"testing"

	"github.com/reoden/go-echo-template/pkg/logger"
	"github.com/reoden/go-echo-template/pkg/mongodb"
)

var MongoDockerTestContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *mongodb.MongoDbOptions, logger logger.Logger) (*mongodb.MongoDbOptions, error) {
		return NewMongoDockerTest().PopulateContainerOptions(ctx, t)
	}
}
