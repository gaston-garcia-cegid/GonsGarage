package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

type CarRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo *PostgresCarRepository
}

func (suite *CarRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(suite.T(), err)

	// Auto-migrate tables
	err = db.AutoMigrate(&CarModel{}, &UserModel{})
	require.NoError(suite.T(), err)

	suite.db = db
	suite.repo = &PostgresCarRepository{db: db}
}

func (suite *CarRepositoryTestSuite) TearDownTest() {
	// Clean up tables after each test
	suite.db.Exec("DELETE FROM cars")
	suite.db.Exec("DELETE FROM users")
}

func (suite *CarRepositoryTestSuite) TestCreate_Success() {
	// Arrange
	ownerID := uuid.New()
	car := &domain.Car{
		ID:           uuid.New(),
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		VIN:          "1HGBH41JXMN109186",
		Color:        "Blue",
		Mileage:      5000,
		OwnerID:      ownerID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Act
	err := suite.repo.Create(context.Background(), car)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify car was created
	var dbCar CarModel
	result := suite.db.First(&dbCar, "id = ?", car.ID)
	assert.NoError(suite.T(), result.Error)
	assert.Equal(suite.T(), car.Make, dbCar.Make)
	assert.Equal(suite.T(), car.Model, dbCar.Model)
	assert.Equal(suite.T(), car.LicensePlate, dbCar.LicensePlate)
}

func (suite *CarRepositoryTestSuite) TestGetByID_Success() {
	// Arrange
	ownerID := uuid.New()
	carID := uuid.New()

	// Create a car in the database
	dbCar := &CarModel{
		ID:           carID,
		Make:         "Honda",
		Model:        "Civic",
		Year:         2022,
		LicensePlate: "XYZ-789",
		Color:        "Red",
		OwnerID:      ownerID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(dbCar)

	// Act
	result, err := suite.repo.GetByID(context.Background(), carID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), carID, result.ID)
	assert.Equal(suite.T(), "Honda", result.Make)
	assert.Equal(suite.T(), "Civic", result.Model)
}

func (suite *CarRepositoryTestSuite) TestGetByID_NotFound() {
	// Arrange
	nonExistentID := uuid.New()

	// Act
	result, err := suite.repo.GetByID(context.Background(), nonExistentID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), domain.ErrCarNotFound, err)
}

func (suite *CarRepositoryTestSuite) TestGetByOwnerID_Success() {
	// Arrange
	ownerID := uuid.New()

	// Create multiple cars for the same owner
	cars := []*CarModel{
		{
			ID:           uuid.New(),
			Make:         "Toyota",
			Model:        "Camry",
			Year:         2023,
			LicensePlate: "ABC-123",
			Color:        "Blue",
			OwnerID:      ownerID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Make:         "Honda",
			Model:        "Civic",
			Year:         2022,
			LicensePlate: "XYZ-789",
			Color:        "Red",
			OwnerID:      ownerID,
			CreatedAt:    time.Now().Add(-time.Hour), // Older car
			UpdatedAt:    time.Now().Add(-time.Hour),
		},
	}

	for _, car := range cars {
		suite.db.Create(car)
	}

	// Act
	result, err := suite.repo.GetByOwnerID(context.Background(), ownerID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	// Should be ordered by created_at DESC (newest first)
	assert.Equal(suite.T(), "Toyota", result[0].Make)
	assert.Equal(suite.T(), "Honda", result[1].Make)
}

func (suite *CarRepositoryTestSuite) TestGetByLicensePlate_Success() {
	// Arrange
	licensePlate := "ABC-123"
	dbCar := &CarModel{
		ID:           uuid.New(),
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: licensePlate,
		Color:        "Blue",
		OwnerID:      uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(dbCar)

	// Act
	result, err := suite.repo.GetByLicensePlate(context.Background(), licensePlate)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), licensePlate, result.LicensePlate)
}

func (suite *CarRepositoryTestSuite) TestGetByLicensePlate_NotFound() {
	// Act
	result, err := suite.repo.GetByLicensePlate(context.Background(), "NON-EXISTENT")

	// Assert
	assert.NoError(suite.T(), err) // Not finding should not be an error
	assert.Nil(suite.T(), result)
}

func (suite *CarRepositoryTestSuite) TestUpdate_Success() {
	// Arrange
	carID := uuid.New()
	originalCar := &CarModel{
		ID:           carID,
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		OwnerID:      uuid.New(),
		CreatedAt:    time.Now().Add(-time.Hour),
		UpdatedAt:    time.Now().Add(-time.Hour),
	}
	suite.db.Create(originalCar)

	updatedCar := &domain.Car{
		ID:           carID,
		Make:         "Toyota",
		Model:        "Camry Hybrid",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Green",
		Mileage:      10000,
		OwnerID:      originalCar.OwnerID,
		CreatedAt:    originalCar.CreatedAt,
		UpdatedAt:    time.Now(),
	}

	// Act
	err := suite.repo.Update(context.Background(), updatedCar)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify update
	var dbCar CarModel
	suite.db.First(&dbCar, "id = ?", carID)
	assert.Equal(suite.T(), "Camry Hybrid", dbCar.Model)
	assert.Equal(suite.T(), "Green", dbCar.Color)
	assert.Equal(suite.T(), 10000, dbCar.Mileage)
}

func (suite *CarRepositoryTestSuite) TestDelete_Success() {
	// Arrange
	carID := uuid.New()
	dbCar := &CarModel{
		ID:           carID,
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		OwnerID:      uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(dbCar)

	// Act
	err := suite.repo.Delete(context.Background(), carID)

	// Assert
	assert.NoError(suite.T(), err)

	// Verify soft delete
	var dbCar2 CarModel
	result := suite.db.First(&dbCar2, "id = ? AND deleted_at IS NULL", carID)
	assert.Error(suite.T(), result.Error)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, result.Error)

	// Verify record still exists with deleted_at set
	result = suite.db.Unscoped().First(&dbCar2, "id = ?", carID)
	assert.NoError(suite.T(), result.Error)
	assert.NotNil(suite.T(), dbCar2.DeletedAt)
}

func TestCarRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CarRepositoryTestSuite))
}
