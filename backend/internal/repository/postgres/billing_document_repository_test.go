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

type BillingDocumentRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo ports.BillingDocumentRepository
}

func (suite *BillingDocumentRepositoryTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), db.AutoMigrate(&domain.BillingDocument{}))
	suite.db = db
	suite.repo = NewPostgresBillingDocumentRepository(db)
}

func (suite *BillingDocumentRepositoryTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM billing_documents")
}

func (suite *BillingDocumentRepositoryTestSuite) TestCreateGetUpdate() {
	id := uuid.New()
	doc := &domain.BillingDocument{
		ID: id, Kind: domain.BillingDocumentKindPayroll, Title: "March", Amount: 0,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(suite.T(), suite.repo.Create(context.Background(), doc))
	got, err := suite.repo.GetByID(context.Background(), id)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), domain.BillingDocumentKindPayroll, got.Kind)
	doc.Title = "April"
	require.NoError(suite.T(), suite.repo.Update(context.Background(), doc))
	got2, err := suite.repo.GetByID(context.Background(), id)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "April", got2.Title)
}

func (suite *BillingDocumentRepositoryTestSuite) TestGetByID_NotFound() {
	_, err := suite.repo.GetByID(context.Background(), uuid.New())
	require.Error(suite.T(), err)
	assert.ErrorIs(suite.T(), err, domain.ErrBillingDocumentNotFound)
}

func (suite *BillingDocumentRepositoryTestSuite) TestDeleteSoft() {
	id := uuid.New()
	doc := &domain.BillingDocument{
		ID: id, Kind: domain.BillingDocumentKindOther, Title: "T", Amount: 1,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	require.NoError(suite.T(), suite.repo.Create(context.Background(), doc))
	require.NoError(suite.T(), suite.repo.Delete(context.Background(), id))
	_, err := suite.repo.GetByID(context.Background(), id)
	require.Error(suite.T(), err)
}

func TestBillingDocumentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(BillingDocumentRepositoryTestSuite))
}
