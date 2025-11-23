package v1

import (
	"context"
	"fmt"

	"github.com/reoden/go-echo-template/catalogs/internal/products/data/datamodels"
	dtoV1 "github.com/reoden/go-echo-template/catalogs/internal/products/dtos/v1"
	"github.com/reoden/go-echo-template/catalogs/internal/products/dtos/v1/fxparams"
	"github.com/reoden/go-echo-template/catalogs/internal/products/features/gettingproductbyid/v1/dtos"
	"github.com/reoden/go-echo-template/catalogs/internal/products/models"
	"github.com/reoden/go-echo-template/pkg/core/cqrs"
	customErrors "github.com/reoden/go-echo-template/pkg/http/httperrors/customerrors"
	"github.com/reoden/go-echo-template/pkg/logger"
	"github.com/reoden/go-echo-template/pkg/mapper"
	"github.com/reoden/go-echo-template/pkg/postgresgorm/gormdbcontext"

	"github.com/mehdihadeli/go-mediatr"
)

type GetProductByIDHandler struct {
	fxparams.ProductHandlerParams
}

func NewGetProductByIDHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*GetProductById, *dtos.GetProductByIdResponseDto] {
	return &GetProductByIDHandler{
		ProductHandlerParams: params,
	}
}

func (c *GetProductByIDHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*GetProductById, *dtos.GetProductByIdResponseDto](
		c,
	)
}

func (c *GetProductByIDHandler) Handle(
	ctx context.Context,
	query *GetProductById,
) (*dtos.GetProductByIdResponseDto, error) {
	product, err := gormdbcontext.FindModelByID[*datamodels.ProductDataModel, *models.Product](
		ctx,
		c.CatalogsDBContext,
		query.ProductID,
	)
	if err != nil {
		return nil, err
	}

	productDto, err := mapper.Map[*dtoV1.ProductDto](product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping product",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"product with id: {%s} fetched",
			query.ProductID,
		),
		logger.Fields{"Id": query.ProductID.String()},
	)

	return &dtos.GetProductByIdResponseDto{Product: productDto}, nil
}
