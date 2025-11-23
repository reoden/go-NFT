package unittest

import (
	"context"
	"testing"

	"github.com/reoden/go-echo-template/config"
	"github.com/reoden/go-echo-template/internal/pkg/config/environment"
	"github.com/reoden/go-echo-template/internal/pkg/logger"
	defaultLogger "github.com/reoden/go-echo-template/internal/pkg/logger/defaultlogger"
	"github.com/reoden/go-echo-template/internal/services/catalogs/products/models"
	"github.com/reoden/go-echo-template/mocks"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UnitTestSharedFixture struct {
	Cfg *config.AppOptions
	Log logger.Logger
	suite.Suite
	ProductRepository *mocks.ProductRepository
	Products          []*models.Product
	Ctx               context.Context
	DB                *gorm.DB
}

func NewUnitTestSharedFixture(t *testing.T) *UnitTestSharedFixture {
	// this fix root working directory problem in our test environment inner our fixture
	environment.FixProjectRootWorkingDirectoryPath()

	log := defaultLogger.GetLogger()
	cfg := &config.AppOptions{Name: "TestApp"}

	unit := &UnitTestSharedFixture{
		Cfg: cfg,
		Log: log,
	}

	return unit
}

// Shared Hooks
func (c *UnitTestSharedFixture) SetupSuite() {
}

func (c *UnitTestSharedFixture) TearDownSuite() {
}

func (c *UnitTestSharedFixture) SetupTest() {
	c.Ctx = context.Background()

	// Initialize any other fresh test data
	c.Products = []*models.Product{}

	// rest mock for each test run for preventing incorrect mock data between tests
	c.ProductRepository = mocks.NewProductRepository(c.T())
}

func (c *UnitTestSharedFixture) TearDownTest() {
}
