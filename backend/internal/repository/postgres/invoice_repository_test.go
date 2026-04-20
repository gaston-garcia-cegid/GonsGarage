package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

type InvoiceRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo ports.InvoiceRepository
}

func (suite *InvoiceRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), db.AutoMigrate(&domain.Invoice{}))
	suite.db = db
	suite.repo = NewPostgresInvoiceRepository(db)
}

func (suite *InvoiceRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM invoices")
}

func (suite *InvoiceRepositoryTestSuite) TestCreateGetListByCustomer_ListForStaff_Delete() {
	cust := uuid.New()
	id := uuid.New()
	inv := &domain.Invoice{
		ID: id, CustomerID: cust, Amount: 99, Status: "open", Notes: "x",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(suite.T(), suite.repo.Create(context.Background(), inv))
	got, err := suite.repo.GetByID(context.Background(), id)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 99.0, got.Amount)
	list, total, err := suite.repo.ListByCustomerID(context.Background(), cust, 10, 0)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), list, 1)
	staffList, staffTotal, err := suite.repo.ListForStaff(context.Background(), 10, 0)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), staffTotal)
	assert.Len(suite.T(), staffList, 1)
	require.NoError(suite.T(), suite.repo.Delete(context.Background(), id))
	_, err = suite.repo.GetByID(context.Background(), id)
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrInvoiceNotFound)
}

func (suite *InvoiceRepositoryTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), uuid.New())
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrInvoiceNotFound)
}

func TestInvoiceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InvoiceRepositoryTestSuite))
}
