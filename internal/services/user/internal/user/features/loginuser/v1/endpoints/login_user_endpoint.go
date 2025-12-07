package endpoints

import (
	"fmt"
	"net/http"

	"github.com/reoden/go-NFT/pkg/core/web/route"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/utils"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/loginuser/v1/commands"
	"github.com/reoden/go-NFT/user/internal/user/features/loginuser/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type loginUserEndpoint struct {
	fxparams.UserRouteParams
}

func NewLoginUserEndpoint(
	params fxparams.UserRouteParams,
) route.Endpoint {
	return &loginUserEndpoint{UserRouteParams: params}
}

func (ep *loginUserEndpoint) MapEndpoint() {
	ep.UserGroup.POST("/login", ep.handler())
}

// LoginUser
// @Tags User
// @Summary user loginuser
// @Description user loginuser check
// @Accept json
// @Produce json
// @Param LoginUserRequestDto body dtos.LoginUserRequestDto true "User data"
// @Success 200 {object} dtos.LoginUserResponseDto
// @Router /api/v1/user/loginuser [post]
func (ep *loginUserEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.LoginUserRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := commands.NewLoginUserWithValidation(
			request.Phone,
			request.Captcha,
		)
		if err != nil {
			return err
		}

		result, err := mediatr.Send[*commands.LoginUser, *dtos.LoginUserResponseDto](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending LoginUser",
			)
		}

		token, err := utils.GenJWTToken(result.UserId)
		if err != nil {
			return customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf("[Login_User_Handler] generate jwt token for userId=%s err=%+v", result.UserId.String(), err),
			)
		}

		c.Response().Header().Set("Authorization", token)

		return c.JSON(http.StatusOK, result)
	}
}
