package gorm

import (
	"context"
	"testing"

	gormPostgres "github.com/reoden/go-NFT/pkg/postgresgorm"
)

var GormDockerTestConatnerOptionsDecorator = func(t *testing.T, ctx context.Context) interface{} {
	return func(c *gormPostgres.GormOptions) (*gormPostgres.GormOptions, error) {
		return NewGormDockerTest().PopulateContainerOptions(ctx, t)
	}
}
