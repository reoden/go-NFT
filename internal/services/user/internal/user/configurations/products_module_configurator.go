package configurations

import (
	fxcontracts "github.com/reoden/go-NFT/pkg/fxapp/contracts"
	grpcServer "github.com/reoden/go-NFT/pkg/grpc"
	"github.com/reoden/go-NFT/user/internal/shared/grpc"
	productsservice "github.com/reoden/go-NFT/user/internal/shared/grpc/genproto"
	"github.com/reoden/go-NFT/user/internal/user/configurations/endpoints"
	"github.com/reoden/go-NFT/user/internal/user/configurations/mappings"
	"github.com/reoden/go-NFT/user/internal/user/configurations/mediator"

	googleGrpc "google.golang.org/grpc"
)

type ProductsModuleConfigurator struct {
	fxcontracts.Application
}

func NewProductsModuleConfigurator(
	fxapp fxcontracts.Application,
) *ProductsModuleConfigurator {
	return &ProductsModuleConfigurator{
		Application: fxapp,
	}
}

func (c *ProductsModuleConfigurator) ConfigureProductsModule() error {
	// config products mappings
	err := mappings.ConfigureProductsMappings()
	if err != nil {
		return err
	}

	// register products request handler on mediator
	c.ResolveFuncWithParamTag(
		mediator.RegisterMediatorHandlers,
		`group:"product-handlers"`,
	)

	return nil
}

func (c *ProductsModuleConfigurator) MapProductsEndpoints() error {
	// config endpoints
	c.ResolveFuncWithParamTag(
		endpoints.RegisterEndpoints,
		`group:"product-routes"`,
	)

	// config Products Grpc Endpoints
	c.ResolveFunc(
		func(userGrpcServer grpcServer.GrpcServer, grpcService *grpc.ProductGrpcServiceServer) error {
			userGrpcServer.GrpcServiceBuilder().
				RegisterRoutes(func(server *googleGrpc.Server) {
					productsservice.RegisterProductsServiceServer(
						server,
						grpcService,
					)
				})

			return nil
		},
	)

	return nil
}
