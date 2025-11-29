package gorm

import (
	"context"
	"testing"

	"github.com/reoden/go-NFT/pkg/logger"
	gormPostgres "github.com/reoden/go-NFT/pkg/postgresgorm"
)

var GormContainerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions, logger logger.Logger) (*gormPostgres.GormOptions, error) {
		return NewGormTestContainers(logger).PopulateContainerOptions(ctx, t)
	}
}
