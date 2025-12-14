package endpoints

import (
	"net/http"

	"github.com/reoden/go-NFT/pkg/core/web/route"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/commands"
	"github.com/reoden/go-NFT/user/internal/user/features/creatinguser/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type createUserEndpoint struct {
	fxparams.UserRouteParams
}

func NewCreateUserEndpoint(
	params fxparams.UserRouteParams,
) route.Endpoint {
	return &createUserEndpoint{UserRouteParams: params}
}

func (ep *createUserEndpoint) MapEndpoint() {
	ep.UserGroup.POST("/register", ep.handler())
}

// CreateUser
// @Tags User
// @Summary Create user
// @Description Create new user item
// @Accept json
// @Produce json
// @Param CreateUserRequestDto body dtos.CreateUserRequestDto true "User data"
// @Success 201 {object} dtos.CreateUserResponseDto
// @Router /api/v1/user/register [post]
func (ep *createUserEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.CreateUserRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := commands.NewCreateUserWithValidation(
			request.Phone,
			request.Captcha,
		)
		if err != nil {
			return err
		}

		result, err := mediatr.Send[*commands.CreateUser, *dtos.CreateUserResponseDto](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending CreateUser",
			)
		}

		return c.JSON(http.StatusCreated, result)
	}
}
