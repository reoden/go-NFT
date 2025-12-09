package endpoints

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/constants"
	"github.com/reoden/go-NFT/pkg/core/web/route"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/pkg/utils"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/finduserbyId/v1/dtos"
	"github.com/reoden/go-NFT/user/internal/user/features/finduserbyId/v1/queries"
)

type findUserByIdEndpoint struct {
	fxparams.UserRouteParams
}

func NewFindUserByIdEndpoint(
	params fxparams.UserRouteParams,
) route.Endpoint {
	return &findUserByIdEndpoint{UserRouteParams: params}
}

func (ep *findUserByIdEndpoint) MapEndpoint() {
	ep.UserGroup.GET("/:user_id", ep.handler())
}

// FindUserById
// @Tags User
// @Summary Find user By UserId
// @Description Find user item by UserId
// @Accept json
// @Produce json
// @Param FindUserByIdRequestDto body dtos.FindUserByIdRequestDto true "User data"
// @Success 201 {object} dtos.FindUserByIdResponseDto
// @Router /api/v1/user/{user_id} [get]
func (ep *findUserByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		_, userId, err := utils.ParseJWTToken(c)
		if err != nil {
			return customErrors.NewUnAuthorizedErrorWrap(
				err,
				constants.ErrJWTTokenInvalid,
			)
		}

		request := &dtos.FindUserByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		if request.UserId != userId {
			return customErrors.NewUnAuthorizedError(
				"token parse user_id does not eq with request.UserId",
			)
		}

		query, err := queries.NewFindUserByIdWithValidation(
			request.UserId,
		)
		if err != nil {
			return err
		}

		result, err := mediatr.Send[*queries.FindUserById, *dtos.FindUserByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending FindUserById",
			)
		}

		return c.JSON(http.StatusOK, result)
	}
}
