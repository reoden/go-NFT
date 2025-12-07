package user

import (
	"github.com/labstack/echo/v4"
	"github.com/reoden/go-NFT/pkg/core/web/route"
	"github.com/reoden/go-NFT/pkg/http/customecho/contracts"
	"github.com/reoden/go-NFT/user/internal/shared/grpc"
	userConstracts "github.com/reoden/go-NFT/user/internal/user/contracts"
	"github.com/reoden/go-NFT/user/internal/user/data/repositories"
	creatingUserV1 "github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/endpoints"
	findUserByIdV1 "github.com/reoden/go-NFT/user/internal/user/features/findUserById/v1/endpoints"
	loginUserV1 "github.com/reoden/go-NFT/user/internal/user/features/loginuser/v1/endpoints"
	sendCaptchaV1 "github.com/reoden/go-NFT/user/internal/user/features/sendcaptcha/v1/endpoints"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"userfx",

	// Other provides
	fx.Provide(repositories.NewPostgresUserRepository),
	fx.Provide(
		fx.Annotate(
			repositories.NewRedisUserRepository,
			fx.As(new(userConstracts.UserCacheRepository)),
		)),
	fx.Provide(grpc.NewUserGrpcService),

	fx.Provide(
		fx.Annotate(func(userServer contracts.EchoHttpServer) *echo.Group {
			var g *echo.Group
			userServer.RouteBuilder().
				RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
					group := v1.Group("/user")
					g = group
				})

			return g
		}, fx.ResultTags(`name:"user-echo-group"`)),
	),

	// add cqrs handlers to DI
	//fx.Provide(
	//	cqrs.AsHandler(
	//		creatinguserv1.NewCreateProductHandler,
	//		"user-handlers",
	//	),
	//	cqrs.AsHandler(
	//		gettingproductsv1.NewGetProductsHandler,
	//		"product-handlers",
	//	),
	//	cqrs.AsHandler(
	//		deletingproductv1.NewDeleteProductHandler,
	//		"product-handlers",
	//	),
	//	cqrs.AsHandler(
	//		gettingproductbyidv1.NewGetProductByIDHandler,
	//		"product-handlers",
	//	),
	//	cqrs.AsHandler(
	//		searchingproductsv1.NewSearchProductsHandler,
	//		"product-handlers",
	//	),
	//	cqrs.AsHandler(
	//		updatingoroductsv1.NewUpdateProductHandler,
	//		"product-handlers",
	//	),
	//),

	// add endpoints to DI
	fx.Provide(
		route.AsRoute(
			creatingUserV1.NewCreateUserEndpoint,
			"user-routes",
		),
		route.AsRoute(
			findUserByIdV1.NewFindUserByIdEndpoint,
			"user-routes",
		),
		route.AsRoute(
			loginUserV1.NewLoginUserEndpoint,
			"user-routes",
		),
		route.AsRoute(
			sendCaptchaV1.NewSendCaptchaEndpoint,
			"user-routes",
		),
		//route.AsRoute(
		//	updatingoroductsv1.NewUpdateProductEndpoint,
		//	"product-routes",
		//),
		//route.AsRoute(
		//	gettingproductsv1.NewGetProductsEndpoint,
		//	"product-routes",
		//),
		//route.AsRoute(
		//	searchingproductsv1.NewSearchProductsEndpoint,
		//	"product-routes",
		//),
		//route.AsRoute(
		//	gettingproductbyidv1.NewGetProductByIdEndpoint,
		//	"product-routes",
		//),
		//route.AsRoute(
		//	deletingproductv1.NewDeleteProductEndpoint,
		//	"product-routes",
		//),
	),
)
