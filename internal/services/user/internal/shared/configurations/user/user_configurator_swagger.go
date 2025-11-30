package user

import (
	customEcho "github.com/reoden/go-NFT/pkg/http/customecho/contracts"
	//"github.com/reoden/go-NFT/user/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (ic *UserServiceConfigurator) configSwagger(routeBuilder *customEcho.RouteBuilder) {
	// https://github.com/swaggo/swag#how-to-use-it-with-gin
	//docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Title = "User Service Api"
	//docs.SwaggerInfo.Description = "User Service Api."

	routeBuilder.RegisterRoutes(func(e *echo.Echo) {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	})
}
