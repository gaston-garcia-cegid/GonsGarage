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

type ReceivedInvoiceRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo ports.ReceivedInvoiceRepository
}

func (suite *ReceivedInvoiceRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), db.AutoMigrate(&domain.ReceivedInvoice{}))
	suite.db = db
	suite.repo = NewPostgresReceivedInvoiceRepository(db)
}

func (suite *ReceivedInvoiceRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM received_invoices")
}

func sampleReceivedInvoice(id uuid.UUID) *domain.ReceivedInvoice {
	return &domain.ReceivedInvoice{
		ID:          id,
		VendorName:  "Paper Co",
		Category:    "supplies",
		Amount:      120.5,
		InvoiceDate: time.Date(2025, 3, 10, 12, 0, 0, 0, time.UTC),
		Notes:       "n",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (suite *ReceivedInvoiceRepositoryTestSuite) TestCreateGetList() {
	id := uuid.New()
	inv := sampleReceivedInvoice(id)
	require.NoError(suite.T(), suite.repo.Create(context.Background(), inv))
	got, err := suite.repo.GetByID(context.Background(), id)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 120.5, got.Amount)
	list, total, err := suite.repo.List(context.Background(), 10, 0)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	assert.Len(suite.T(), list, 1)
}

func (suite *ReceivedInvoiceRepositoryTestSuite) TestDeleteSoftExcludesFromList() {
	id := uuid.New()
	require.NoError(suite.T(), suite.repo.Create(context.Background(), sampleReceivedInvoice(id)))
	require.NoError(suite.T(), suite.repo.Delete(context.Background(), id))
	_, err := suite.repo.GetByID(context.Background(), id)
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrReceivedInvoiceNotFound)
	list, total, err := suite.repo.List(context.Background(), 10, 0)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(0), total)
	assert.Len(suite.T(), list, 0)
}

func (suite *ReceivedInvoiceRepositoryTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), uuid.New())
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrReceivedInvoiceNotFound)
}

func TestReceivedInvoiceRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ReceivedInvoiceRepositoryTestSuite))
}
