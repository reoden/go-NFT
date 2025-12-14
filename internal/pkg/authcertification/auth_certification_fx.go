package authcertification

import (
    "context"
    "fmt"

    "github.com/reoden/go-NFT/pkg/http/client"
    "github.com/reoden/go-NFT/pkg/logger"
    "go.uber.org/fx"
    "go.uber.org/zap"
)

var (
	Module = fx.Module(
		"authcertificationfx",
		authProviders,
		authInvokes,
	)

	authProviders = fx.Provide(
		provideConfig,
		client.NewHttpClient,
		NewAuthCertificationService,
	)

	authInvokes = fx.Invoke(registerHooks)
)

func registerHooks(
	lc fx.Lifecycle,
	authServer AuthCertificationService,
	logger logger.Logger,
) {
	implName := "unknown"
	switch impl := authServer.(type) {
	case *MockAuthCertificationServiceImpl:
		implName = "MockAuthCertificationServiceImpl"
		logger.Info("using MockAuthCertificationServiceImpl")
	case *AuthCertificationServiceImpl:
		implName = "AuthCertificationServiceImpl"
		logger.Info("using AuthCertificationServiceImpl")
	default:
		logger.Warn("unknown AuthCertificationService implementation", zap.String("type", fmt.Sprintf("%T", impl)))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Infof("successfully register AutService = '%s'", implName)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Infof("successfully unregister AutService = '%s'", implName)

			return nil
		},
	})
}