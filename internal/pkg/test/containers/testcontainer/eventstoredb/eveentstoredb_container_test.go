package eventstoredb

import (
	"context"
	"testing"

	"github.com/reoden/go-echo-template/pkg/config"
	"github.com/reoden/go-echo-template/pkg/config/environment"
	"github.com/reoden/go-echo-template/pkg/core"
	"github.com/reoden/go-echo-template/pkg/eventstroredb"
	"github.com/reoden/go-echo-template/pkg/logger/external/fxlog"
	"github.com/reoden/go-echo-template/pkg/logger/zap"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func Test_Custom_EventStoreDB_Container(t *testing.T) {
	var esdbClient *esdb.Client
	ctx := context.Background()

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		eventstroredb.ModuleFunc(func() {
		}),
		fx.Decorate(EventstoreDBContainerOptionsDecorator(t, ctx)),
		fx.Populate(&esdbClient),
	).RequireStart()
}
