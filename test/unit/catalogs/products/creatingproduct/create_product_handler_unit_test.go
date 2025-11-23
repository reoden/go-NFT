//go:build unit
// +build unit

package creatingproduct

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	customErrors "github.com/reoden/go-echo-template/internal/pkg/http/httperrors/customerrors"
	commands2 "github.com/reoden/go-echo-template/internal/services/catalogs/products/features/creatingproduct/commands"
	"github.com/reoden/go-echo-template/internal/services/catalogs/products/features/creatingproduct/dtos"
	"github.com/reoden/go-echo-template/test/testfixtures/unittest"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type createProductHandlerUnitTests struct {
	*unittest.UnitTestSharedFixture
	handler mediatr.RequestHandler[*commands2.CreateProductCommand, *dtos.CreateProductCommandResponse]
}

func TestCreateProductHandlerUnit(t *testing.T) {
	suite.Run(t, &createProductHandlerUnitTests{
		UnitTestSharedFixture: unittest.NewUnitTestSharedFixture(t),
	},
	)
}

func (c *createProductHandlerUnitTests) SetupTest() {
	// call base SetupTest hook before running child hook
	c.UnitTestSharedFixture.SetupTest()
	c.handler = commands2.NewCreateProductCommandHandler(c.ProductRepository)
}

func (c *createProductHandlerUnitTests) TearDownTest() {
	// call base TearDownTest hook before running child hook
	c.UnitTestSharedFixture.TearDownTest()
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Create_New_Product_With_Valid_Data() {
	id := uuid.NewV4()

	createProduct := &commands2.CreateProductCommand{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := commands2.MapCreateProductToProduct(createProduct)

	c.ProductRepository.On("CreateProduct", c.Ctx, product).
		Return(product, nil).
		Once()

	res, err := c.handler.Handle(c.Ctx, createProduct)

	c.Require().NoError(err)

	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
	c.ProductRepository.AssertCalled(c.T(), "CreateProduct", c.Ctx, product)

	c.Require().NoError(err)
	c.Require().NotNil(res)
	c.Equal(product.ProductID, res.ProductID)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_When_Repository_Fails() {
	id := uuid.NewV4()
	createProduct := &commands2.CreateProductCommand{
		ProductID:   id,
		Name:        gofakeit.Name(),
		CreatedAt:   time.Now(),
		Description: gofakeit.EmojiDescription(),
		Price:       gofakeit.Price(100, 1000),
	}

	product := commands2.MapCreateProductToProduct(createProduct)

	c.ProductRepository.On("CreateProduct", c.Ctx, product).
		Return(nil, customErrors.NewNotFoundError(fmt.Sprintf("product with id %s not found", createProduct.ProductID))).
		Once()

	res, err := c.handler.Handle(c.Ctx, createProduct)

	c.Require().Error(err)
	c.Require().ErrorContains(err, fmt.Sprintf("product with id %s not found", createProduct.ProductID))
	c.True(customErrors.IsApplicationError(err, http.StatusInternalServerError))
	c.Require().Nil(res)
	c.ProductRepository.AssertCalled(c.T(), "CreateProduct", c.Ctx, product)
	c.ProductRepository.AssertNumberOfCalls(c.T(), "CreateProduct", 1)
}

func (c *createProductHandlerUnitTests) Test_Handle_Should_Return_Error_When_Command_Is_Nil() {
	res, err := c.handler.Handle(c.Ctx, nil)

	c.Require().Error(err)
	c.Require().ErrorContains(err, "command cannot be nil")
	c.True(customErrors.IsApplicationError(err, http.StatusBadRequest))
	c.Require().Nil(res)
}
