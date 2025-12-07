package endpoints

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/reoden/go-NFT/pkg/core/web/route"
	customErrors "github.com/reoden/go-NFT/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-NFT/user/internal/user/dtos/v1/fxparams"
	"github.com/reoden/go-NFT/user/internal/user/features/sendcaptcha/v1/commands"
	"github.com/reoden/go-NFT/user/internal/user/features/sendcaptcha/v1/dtos"
)

type sendCaptchaEndpoint struct {
	fxparams.UserRouteParams
}

func NewSendCaptchaEndpoint(
	params fxparams.UserRouteParams,
) route.Endpoint {
	return &sendCaptchaEndpoint{UserRouteParams: params}
}

func (ep *sendCaptchaEndpoint) MapEndpoint() {
	ep.UserGroup.POST("/captcha", ep.handler())
}

// SendCaptcha
// @Tags User
// @Summary send user captcha
// @Description send user captcha
// @Accept json
// @Produce json
// @Param SendCaptchaRequestDto body dtos.SendCaptchaRequestDto true "User data"
// @Success 200 {object} dtos.SendCaptchaResponseDto
// @Router /api/v1/user/captcha [post]
func (ep *sendCaptchaEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.SendCaptchaRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := commands.NewSendCaptchaWithValidation(
			request.Phone,
		)
		if err != nil {
			return err
		}

		result, err := mediatr.Send[*commands.SendCaptcha, *dtos.SendCaptchaResponseDto](
			ctx,
			command,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending SendCaptcha",
			)
		}

		return c.JSON(http.StatusOK, result)
	}
}
