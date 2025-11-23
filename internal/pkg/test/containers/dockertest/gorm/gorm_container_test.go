package gorm

import (
	"context"
	"testing"

	"github.com/reoden/go-echo-template/pkg/config"
	"github.com/reoden/go-echo-template/pkg/config/environment"
	"github.com/reoden/go-echo-template/pkg/core"
	"github.com/reoden/go-echo-template/pkg/logger/external/fxlog"
	"github.com/reoden/go-echo-template/pkg/logger/zap"
	gormPostgres "github.com/reoden/go-echo-template/pkg/postgresgorm"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

func Test_Gorm_Container(t *testing.T) {
	ctx := context.Background()
	var gorm *gorm.DB

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		gormPostgres.Module,
		fx.Decorate(GormDockerTestConatnerOptionsDecorator(t, ctx)),
		fx.Populate(&gorm),
	).RequireStart()

	assert.NotNil(t, gorm)
}
