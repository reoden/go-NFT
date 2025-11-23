package v1

import (
	"net/http"

	"github.com/reoden/go-echo-template/catalogs/internal/products/dtos/v1/fxparams"
	"github.com/reoden/go-echo-template/catalogs/internal/products/features/gettingproductbyid/v1/dtos"
	"github.com/reoden/go-echo-template/pkg/core/web/route"
	customErrors "github.com/reoden/go-echo-template/pkg/http/httperrors/customerrors"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getProductByIdEndpoint struct {
	fxparams.ProductRouteParams
}

func NewGetProductByIdEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &getProductByIdEndpoint{ProductRouteParams: params}
}

func (ep *getProductByIdEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product by id
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		query, err := NewGetProductByIdWithValidation(request.ProductId)
		if err != nil {
			return err
		}

		queryResult, err := mediatr.Send[*GetProductById, *dtos.GetProductByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProductById",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
