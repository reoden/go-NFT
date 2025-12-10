package configurations

import (
	"github.com/reoden/go-NFT/pkg/bloom"
	fxcontracts "github.com/reoden/go-NFT/pkg/fxapp/contracts"
	grpcServer "github.com/reoden/go-NFT/pkg/grpc"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/otel/tracing"
	"github.com/reoden/go-NFT/user/internal/shared/data/dbcontext"
	"github.com/reoden/go-NFT/user/internal/shared/grpc"
	userservice "github.com/reoden/go-NFT/user/internal/shared/grpc/genproto"
	"github.com/reoden/go-NFT/user/internal/user/configurations/endpoints"
	"github.com/reoden/go-NFT/user/internal/user/configurations/mappings"
	"github.com/reoden/go-NFT/user/internal/user/configurations/mediator"
	"github.com/reoden/go-NFT/user/internal/user/contracts"

	googleGrpc "google.golang.org/grpc"
)

type UserModuleConfigurator struct {
	fxcontracts.Application
}

func NewUserModuleConfigurator(
	fxapp fxcontracts.Application,
) *UserModuleConfigurator {
	return &UserModuleConfigurator{
		Application: fxapp,
	}
}

func (c *UserModuleConfigurator) ConfigureUserModule() {
	c.ResolveFunc(
		func(logger logger.Logger,
			userDBContext *dbcontext.UserGormDBContext,
			userRepository contracts.UserRepository,
			userOperationRepository contracts.UserOperateStreamRepository,
			cacheRepository contracts.UserCacheRepository,
			bloomFilter *bloom.BloomFilterFactory,
			tracer tracing.AppTracer,
		) error {
			// config User Mediators
			err := mediator.ConfigUserMediator(logger, userDBContext, userRepository, userOperationRepository, cacheRepository, bloomFilter, tracer)
			if err != nil {
				return err
			}

			// config user mappings
			err = mappings.ConfigureUserMappings()
			if err != nil {
				return err
			}

			return nil
		})
}

func (c *UserModuleConfigurator) MapUserEndpoints() {
	// config endpoints
	c.ResolveFuncWithParamTag(
		endpoints.RegisterEndpoints,
		`group:"user-routes"`,
	)

	// config User Grpc Endpoints
	c.ResolveFunc(
		func(userGrpcServer grpcServer.GrpcServer, grpcService *grpc.UserGrpcServiceServer) error {
			userGrpcServer.GrpcServiceBuilder().
				RegisterRoutes(func(server *googleGrpc.Server) {
					userservice.RegisterUserServiceServer(
						server,
						grpcService,
					)
				})

			return nil
		},
	)
}
