package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/core/web/route"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/utils"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/logout/v1/commands"
	"github.com/reoden/go-NFT/user/internal/user/features/logout/v1/dtos"
)

type logoutUserEndpoint struct {
	fxparams.UserRouteParams
}

func NewLogoutUserEndpoint(
	params fxparams.UserRouteParams,
) route.Endpoint {
	return &logoutUserEndpoint{UserRouteParams: params}
}

func (ep *logoutUserEndpoint) MapEndpoint() {
	ep.UserGroup.POST("/logout", ep.handler())
}

// LogoutUser
// @Tags User
// @Summary user logoutuser
// @Description user logoutuser check
// @Accept json
// @Produce json
// @Param LogoutUserRequestDto body dtos.LogoutUserRequestDto true "User data"
// @Success 200 {object} dtos.LogoutUserResponseDto
// @Router /api/v1/user/logout [post]
func (ep *logoutUserEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		token, userId, err := utils.ParseJWTToken(c)
		if err != nil {
			return customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf("[Logout_User_Handler] parse jwt token err=%+v", err),
			)
		}

		command, err := commands.NewLogoutUserWithValidation(token, userId)
		if err != nil {
			return err
		}

		result, err := mediatr.Send[*commands.LogoutUser, *dtos.LogoutUserResponseDto](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending LogoutUser",
			)
		}

		return c.JSON(http.StatusOK, result)
	}
}
