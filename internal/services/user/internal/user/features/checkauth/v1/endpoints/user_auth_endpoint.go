package endpoints

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/constants"
	"github.com/reoden/go-NFT/pkg/core/web/route"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/logger"
	"github.com/reoden/go-NFT/pkg/utils"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/checkauth/v1/commands"
	"github.com/reoden/go-NFT/user/internal/user/features/checkauth/v1/dtos"
)

type authEndpoint struct {
	fxparams.UserRouteParams
}

func NewAuthEndpoint(
	params fxparams.UserRouteParams,
) route.Endpoint {
	return &authEndpoint{UserRouteParams: params}
}

func (ep *authEndpoint) MapEndpoint() {
	ep.UserGroup.POST("/auth", ep.handler())
}

// Auth
// @Tags User
// @Summary user Auth certification
// @Description aut certification
// @Accept json
// @Produce json
// @Param AuthRequestDto body dtos.AuthRequestDto true "real_name and id_card_number"
// @Success 201 {object} dtos.AuthResponseDto
// @Security BearerAuth
// @Router /api/v1/user/auth [post]
func (ep *authEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		_, userId, err := utils.ParseJWTToken(c)
		if err != nil {
			return customErrors.NewUnAuthorizedErrorWrap(
				err,
				constants.ErrJWTTokenInvalid,
			)
		}

		ep.Logger.Infow("user auth", logger.Fields{"userId": userId})

		request := &dtos.AuthRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := commands.NewAuthUserWithValidation(
			request.RealName,
			request.IdCardNo,
			userId,
		)

		if err != nil {
			return err
		}

		result, err := mediatr.Send[*commands.AuthUser, *dtos.AuthResponseDto](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending AuthUser",
			)
		}

		return c.JSON(http.StatusOK, result)
	}
}
