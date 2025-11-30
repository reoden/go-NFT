package app

import (
	"context"

	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/shared/configurations/user"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	// configure dependencies
	appBuilder := NewUserApplicationBuilder()
	appBuilder.ProvideModule(user.UserServiceModule)

	app := appBuilder.Build()

	// configure application
	app.ConfigureUser()

	app.MapUserEndpoints()

	app.Logger().Info("Starting user_service application")
	app.ResolveFunc(func(tracer tracing.AppTracer) {
		_, span := tracer.Start(context.Background(), "Application started")
		span.End()
	})

	app.Run()
}
