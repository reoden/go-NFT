package user

import (
    "fmt"
    "net/http"
    "strings"

    "github.com/reoden/go-NFT/pkg/fxapp/contracts"
    echocontracts "github.com/reoden/go-NFT/pkg/http/customecho/contracts"
    "github.com/reoden/go-NFT/pkg/http/customecho/middlewares/auth"
    migrationcontracts "github.com/reoden/go-NFT/pkg/migration/contracts"
    "github.com/reoden/go-NFT/user/config"
    "github.com/reoden/go-NFT/user/internal/shared/configurations/user/infrastructure"
    "github.com/reoden/go-NFT/user/internal/user/configurations"
    "gorm.io/gorm"

    "github.com/labstack/echo/v4"
)

type UserServiceConfigurator struct {
    contracts.Application
    infrastructureConfigurator *infrastructure.InfrastructureConfigurator
    userModuleConfigurator     *configurations.UserModuleConfigurator
}

func NewUserServiceConfigurator(
    app contracts.Application,
) *UserServiceConfigurator {
    infraConfigurator := infrastructure.NewInfrastructureConfigurator(app)
    productModuleConfigurator := configurations.NewProductsModuleConfigurator(
        app,
    )

    return &UserServiceConfigurator{
        Application:                app,
        infrastructureConfigurator: infraConfigurator,
        userModuleConfigurator:     productModuleConfigurator,
    }
}

func (ic *UserServiceConfigurator) ConfigureUser() {
    // Shared
    // Infrastructure
    ic.infrastructureConfigurator.ConfigInfrastructures()

    // Shared
    // Catalogs configurations
    ic.ResolveFunc(
        func(db *gorm.DB, postgresMigrationRunner migrationcontracts.PostgresMigrationRunner) error {
            err := ic.migrateUser(postgresMigrationRunner)
            if err != nil {
                return err
            }

            //if ic.Environment() != environment.Test {
            //	err = ic.seedUser(db)
            //	if err != nil {
            //		return err
            //	}
            //}

            return nil
        },
    )

    // Modules
    // Product module
    ic.userModuleConfigurator.ConfigureUserModule()
}

func (ic *UserServiceConfigurator) MapUserEndpoints() {
    // Shared
    ic.ResolveFunc(
        func(userServer echocontracts.EchoHttpServer, options *config.AppOptions) error {
            userServer.SetupDefaultMiddlewares()
            authSkipper := func(c echo.Context) bool {
                return func(ec echo.Context) bool {
                    path := ec.Request().URL.Path
                    method := ec.Request().Method
                    if strings.HasPrefix(path, "/api/v1/user") && method == echo.POST {
                        return true
                    }
                    return false
                }(c)
            }
            userServer.AddMiddlewares(auth.EchoAuth(authSkipper))

            // config user root endpoint
            userServer.RouteBuilder().
                RegisterRoutes(func(e *echo.Echo) {
                    e.GET("", func(ec echo.Context) error {
                        return ec.String(
                            http.StatusOK,
                            fmt.Sprintf(
                                "%s is running...",
                                options.GetMicroserviceNameUpper(),
                            ),
                        )
                    })
                })

            // config user swagger
            ic.configSwagger(userServer.RouteBuilder())

            return nil
        },
    )

    // Modules
    // Products CatalogsServiceModule endpoints
    ic.userModuleConfigurator.MapUserEndpoints()
}
