package configurations

import (
	"github.com/reoden/go-echo-template/catalogs/internal/products/configurations/endpoints"
	"github.com/reoden/go-echo-template/catalogs/internal/products/configurations/mappings"
	"github.com/reoden/go-echo-template/catalogs/internal/products/configurations/mediator"
	"github.com/reoden/go-echo-template/catalogs/internal/shared/grpc"
	productsservice "github.com/reoden/go-echo-template/catalogs/internal/shared/grpc/genproto"
	fxcontracts "github.com/reoden/go-echo-template/pkg/fxapp/contracts"
	grpcServer "github.com/reoden/go-echo-template/pkg/grpc"

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
		func(catalogsGrpcServer grpcServer.GrpcServer, grpcService *grpc.ProductGrpcServiceServer) error {
			catalogsGrpcServer.GrpcServiceBuilder().
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
