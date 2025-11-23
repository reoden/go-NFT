package mongo

import (
	"context"
	"testing"

	"github.com/reoden/go-echo-template/pkg/config"
	"github.com/reoden/go-echo-template/pkg/config/environment"
	"github.com/reoden/go-echo-template/pkg/core"
	"github.com/reoden/go-echo-template/pkg/logger/external/fxlog"
	"github.com/reoden/go-echo-template/pkg/logger/zap"
	"github.com/reoden/go-echo-template/pkg/mongodb"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func Test_Mongo_Container(t *testing.T) {
	ctx := context.Background()
	var mongoClient *mongo.Client

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		mongodb.Module,
		fx.Decorate(MongoDockerTestContainerOptionsDecorator(t, ctx)),
		fx.Populate(&mongoClient),
	).RequireStart()

	assert.NotNil(t, mongoClient)
}
