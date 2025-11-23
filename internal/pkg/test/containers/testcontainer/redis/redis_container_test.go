package redis

import (
	"context"
	"testing"

	"github.com/reoden/go-echo-template/pkg/config"
	"github.com/reoden/go-echo-template/pkg/config/environment"
	"github.com/reoden/go-echo-template/pkg/core"
	"github.com/reoden/go-echo-template/pkg/logger/external/fxlog"
	"github.com/reoden/go-echo-template/pkg/logger/zap"
	redis2 "github.com/reoden/go-echo-template/pkg/redis"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func Test_Custom_Redis_Container(t *testing.T) {
	ctx := context.Background()
	var redisClient redis.UniversalClient

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		redis2.Module,
		fx.Decorate(RedisContainerOptionsDecorator(t, ctx)),
		fx.Populate(&redisClient),
	).RequireStart()

	assert.NotNil(t, redisClient)
}
