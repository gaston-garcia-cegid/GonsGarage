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

type PartItemRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo ports.PartItemRepository
}

func (suite *PartItemRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), db.AutoMigrate(&domain.PartItem{}))
	suite.db = db
	suite.repo = NewPostgresPartItemRepository(db)
}

func (suite *PartItemRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM part_items")
}

func newTestPartItem() *domain.PartItem {
	id := uuid.New()
	return &domain.PartItem{
		ID:        id,
		Reference: "REF-1",
		Brand:     "ACME",
		Name:      "Spark plug",
		Barcode:   "5900000000001",
		Quantity:  4,
		UOM:       domain.PartUOMUnit,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (suite *PartItemRepositoryTestSuite) TestCreateAndGetByID() {
	p := newTestPartItem()
	require.NoError(suite.T(), suite.repo.Create(context.Background(), p))
	got, err := suite.repo.GetByID(context.Background(), p.ID)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), got)
	assert.Equal(suite.T(), "REF-1", got.Reference)
	assert.Equal(suite.T(), 4.0, got.Quantity)
}

func (suite *PartItemRepositoryTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), uuid.New())
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrPartItemNotFound)
}

func (suite *PartItemRepositoryTestSuite) TestGetByBarcode() {
	p := newTestPartItem()
	p.Barcode = "  BC99  "
	require.NoError(suite.T(), suite.repo.Create(context.Background(), p))
	got, err := suite.repo.GetByBarcode(context.Background(), "BC99")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), p.ID, got.ID)
}

func (suite *PartItemRepositoryTestSuite) TestGetByBarcode_EmptyNotFound() {
	_, err := suite.repo.GetByBarcode(context.Background(), "   ")
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrPartItemNotFound)
}

func (suite *PartItemRepositoryTestSuite) TestUpdate_NotFound() {
	p := newTestPartItem()
	p.ID = uuid.New()
	err := suite.repo.Update(context.Background(), p)
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrPartItemNotFound)
}

func (suite *PartItemRepositoryTestSuite) TestDelete_SoftThenInvisible() {
	p := newTestPartItem()
	require.NoError(suite.T(), suite.repo.Create(context.Background(), p))
	require.NoError(suite.T(), suite.repo.Delete(context.Background(), p.ID))
	_, err := suite.repo.GetByID(context.Background(), p.ID)
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrPartItemNotFound)
}

func (suite *PartItemRepositoryTestSuite) TestList_BarcodeAndSearchFilters() {
	a := newTestPartItem()
	a.Reference = "NGK-1"
	a.Name = "Bujía"
	a.Barcode = "111"
	require.NoError(suite.T(), suite.repo.Create(context.Background(), a))

	b := newTestPartItem()
	b.ID = uuid.New()
	b.Reference = "OIL-1"
	b.Name = "Aceite 5W30"
	b.Barcode = "222"
	b.Brand = "Castrol"
	require.NoError(suite.T(), suite.repo.Create(context.Background(), b))

	bc := "111"
	list, total, err := suite.repo.List(context.Background(), ports.PartItemListFilters{Barcode: &bc, Limit: 10, Offset: 0})
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total)
	require.Len(suite.T(), list, 1)
	assert.Equal(suite.T(), "111", list[0].Barcode)

	search := "Aceite"
	list2, total2, err := suite.repo.List(context.Background(), ports.PartItemListFilters{Search: &search, Limit: 10, Offset: 0})
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(1), total2)
	require.Len(suite.T(), list2, 1)
	assert.Equal(suite.T(), "OIL-1", list2[0].Reference)
}

func TestPartItemRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PartItemRepositoryTestSuite))
}
