package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/handlers"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/postgres"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/services/car"
)

type CarIntegrationTestSuite struct {
	suite.Suite
	app    *gin.Engine
	db     *gorm.DB
	client *http.Client
}

func (suite *CarIntegrationTestSuite) SetupSuite() {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: nil, // Disable GORM logging for tests
	})
	require.NoError(suite.T(), err)

	// Auto-migrate
	err = db.AutoMigrate(&postgres.CarModel{}, &postgres.UserModel{})
	require.NoError(suite.T(), err)

	// Setup repositories
	carRepo := postgres.NewPostgresCarRepository(db)
	userRepo := postgres.NewPostgresUserRepository(db)

	// ✅ Fixed: Create car service following Agent.md Clean Architecture
	carService := car.NewCarService(carRepo, userRepo, nil)

	// Setup handlers
	carHandler := handlers.NewCarHandler(carService)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add test middleware to set user context
	router.Use(suite.testAuthMiddleware())

	// Setup routes (following Agent.md REST conventions)
	v1 := router.Group("/api/v1")
	cars := v1.Group("/cars")
	{
		cars.POST("", ginHandler(carHandler.CreateCar))
		cars.GET("", ginHandler(carHandler.ListCars))
		cars.GET("/:id", ginHandler(carHandler.GetCar))
		cars.PUT("/:id", ginHandler(carHandler.UpdateCar))
		cars.DELETE("/:id", ginHandler(carHandler.DeleteCar))
	}

	suite.app = router
	suite.db = db
	suite.client = &http.Client{Timeout: 10 * time.Second}
}

func (suite *CarIntegrationTestSuite) TearDownTest() {
	// Clean up test data between tests
	suite.db.Exec("DELETE FROM cars")
	suite.db.Exec("DELETE FROM users")
}

func (suite *CarIntegrationTestSuite) testAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add test user to context (following Agent.md context handling)
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			parsedID, err := uuid.Parse(userID)
			if err == nil {
				c.Set("userID", parsedID)
			}
		}
		c.Next()
	}
}

// ✅ Test follows Agent.md TDD rules: Arrange-Act-Assert pattern
func (suite *CarIntegrationTestSuite) TestCreateCar_Success() {
	// Arrange
	clientID := uuid.New()

	// Create test user (following Agent.md naming conventions)
	testUser := &postgres.UserModel{
		ID:        clientID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(testUser)

	// ✅ Fixed: camelCase naming as per Agent.md
	carData := map[string]interface{}{
		"make":         "Toyota",
		"model":        "Camry",
		"year":         2023,
		"licensePlate": "ABC-123", // camelCase
		"color":        "Blue",
		"mileage":      0,
	}

	jsonData, _ := json.Marshal(carData)

	// Act
	req, _ := http.NewRequest("POST", "/api/v1/cars", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", clientID.String())

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Toyota", response["make"])
	assert.Equal(suite.T(), "Camry", response["model"])
	assert.Equal(suite.T(), float64(2023), response["year"])
	assert.Equal(suite.T(), "ABC-123", response["licensePlate"]) // camelCase response
	assert.Equal(suite.T(), "Blue", response["color"])
	assert.NotEmpty(suite.T(), response["id"])
	assert.Equal(suite.T(), clientID.String(), response["ownerID"]) // camelCase
}

func (suite *CarIntegrationTestSuite) TestCreateCar_Fail_InvalidData() {
	// Arrange
	clientID := uuid.New()

	// Create test user
	testUser := &postgres.UserModel{
		ID:        clientID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(testUser)

	// Invalid car data (missing required fields)
	carData := map[string]interface{}{
		"make": "Toyota",
		// Missing model, year, licensePlate, color
	}

	jsonData, _ := json.Marshal(carData)

	// Act
	req, _ := http.NewRequest("POST", "/api/v1/cars", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", clientID.String())

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *CarIntegrationTestSuite) TestGetCar_Success_ClientAccessesOwnCar() {
	// Arrange
	clientID := uuid.New()
	carID := uuid.New()

	// Create test user
	testUser := &postgres.UserModel{
		ID:        clientID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(testUser)

	// Create test car
	testCar := &postgres.CarModel{
		ID:           carID,
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		OwnerID:      clientID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(testCar)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/cars/"+carID.String(), nil)
	req.Header.Set("X-User-ID", clientID.String())

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), carID.String(), response["id"])
	assert.Equal(suite.T(), "Toyota", response["make"])
	assert.Equal(suite.T(), clientID.String(), response["ownerID"]) // camelCase
}

// ✅ Test verifies client permission rules from Agent.md
func (suite *CarIntegrationTestSuite) TestGetCar_Fail_ClientAccessesOtherClientCar() {
	// Arrange
	clientID := uuid.New()
	otherClientID := uuid.New()
	carID := uuid.New()

	// Create test users
	testUser := &postgres.UserModel{
		ID:        clientID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(testUser)

	otherUser := &postgres.UserModel{
		ID:        otherClientID,
		Email:     "other@example.com",
		FirstName: "Jane",
		LastName:  "Smith",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(otherUser)

	// Create car owned by other client
	testCar := &postgres.CarModel{
		ID:           carID,
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		OwnerID:      otherClientID, // Different owner
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(testCar)

	// Act
	req, _ := http.NewRequest("GET", "/api/v1/cars/"+carID.String(), nil)
	req.Header.Set("X-User-ID", clientID.String()) // Request as first client

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	// Assert - Should be forbidden (client can only access own cars)
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func (suite *CarIntegrationTestSuite) TestListCars_ClientOnlySeesOwnCars() {
	// Arrange
	clientID := uuid.New()
	otherClientID := uuid.New()

	// Create test users
	testUser := &postgres.UserModel{
		ID:        clientID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(testUser)

	otherUser := &postgres.UserModel{
		ID:        otherClientID,
		Email:     "other@example.com",
		FirstName: "Jane",
		LastName:  "Smith",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(otherUser)

	// Create cars for both clients
	clientCar := &postgres.CarModel{
		ID:           uuid.New(),
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		OwnerID:      clientID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(clientCar)

	otherCar := &postgres.CarModel{
		ID:           uuid.New(),
		Make:         "Honda",
		Model:        "Civic",
		Year:         2022,
		LicensePlate: "XYZ-789",
		Color:        "Red",
		OwnerID:      otherClientID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(otherCar)

	// Act - Request as first client
	req, _ := http.NewRequest("GET", "/api/v1/cars", nil)
	req.Header.Set("X-User-ID", clientID.String())

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	// Client should only see their own car (privacy rule from Agent.md)
	assert.Len(suite.T(), response, 1)
	assert.Equal(suite.T(), "Toyota", response[0]["make"])
	assert.Equal(suite.T(), clientID.String(), response[0]["ownerID"]) // camelCase
}

func (suite *CarIntegrationTestSuite) TestUpdateCar_Success_ClientUpdatesOwnCar() {
	// Arrange
	clientID := uuid.New()
	carID := uuid.New()

	// Create test user
	testUser := &postgres.UserModel{
		ID:        clientID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(testUser)

	// Create test car
	testCar := &postgres.CarModel{
		ID:           carID,
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		Mileage:      1000,
		OwnerID:      clientID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(testCar)

	// Updated car data
	updatedCarData := map[string]interface{}{
		"make":         "Toyota",
		"model":        "Camry",
		"year":         2023,
		"licensePlate": "ABC-123",
		"color":        "Red", // Changed color
		"mileage":      2000,  // Updated mileage
	}

	jsonData, _ := json.Marshal(updatedCarData)

	// Act
	req, _ := http.NewRequest("PUT", "/api/v1/cars/"+carID.String(), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", clientID.String())

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "Red", response["color"])           // Updated
	assert.Equal(suite.T(), float64(2000), response["mileage"]) // Updated
	assert.Equal(suite.T(), carID.String(), response["id"])     // Preserved
}

func (suite *CarIntegrationTestSuite) TestDeleteCar_Success_ClientDeletesOwnCar() {
	// Arrange
	clientID := uuid.New()
	carID := uuid.New()

	// Create test user
	testUser := &postgres.UserModel{
		ID:        clientID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Role:      "client",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.db.Create(testUser)

	// Create test car
	testCar := &postgres.CarModel{
		ID:           carID,
		Make:         "Toyota",
		Model:        "Camry",
		Year:         2023,
		LicensePlate: "ABC-123",
		Color:        "Blue",
		OwnerID:      clientID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	suite.db.Create(testCar)

	// Act
	req, _ := http.NewRequest("DELETE", "/api/v1/cars/"+carID.String(), nil)
	req.Header.Set("X-User-ID", clientID.String())

	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusNoContent, w.Code)

	// Verify car is soft deleted
	var deletedCar postgres.CarModel
	result := suite.db.Unscoped().Where("id = ?", carID).First(&deletedCar)
	assert.NoError(suite.T(), result.Error)
	assert.NotNil(suite.T(), deletedCar.DeletedAt) // Should be soft deleted
}

// Helper function to convert handler to gin.HandlerFunc (following Agent.md patterns)
func ginHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c.Writer, c.Request)
	}
}

// ✅ Test suite runner following Agent.md conventions
func TestCarIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(CarIntegrationTestSuite))
}
