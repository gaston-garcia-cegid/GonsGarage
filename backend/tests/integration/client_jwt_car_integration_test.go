//go:build cgo

package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/handler"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/middleware"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/repository/postgres"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/service/car"
)

const jwtTestSecret = "integration-jwt-secret-test-only"

func bearerToken(userID uuid.UUID, role string) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID.String(),
		"role":   role,
	})
	s, err := t.SignedString([]byte(jwtTestSecret))
	if err != nil {
		panic(err)
	}
	return "Bearer " + s
}

type ClientJWTCarIntegrationSuite struct {
	suite.Suite
	app *gin.Engine
	db  *gorm.DB
}

func (suite *ClientJWTCarIntegrationSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: nil})
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), db.AutoMigrate(&postgres.CarModel{}, &postgres.UserModel{}))

	carRepo := postgres.NewPostgresCarRepository(db)
	userRepo := postgres.NewPostgresUserRepository(db)
	carSvc := car.NewCarService(carRepo, userRepo, nil)
	carHandler := handler.NewCarHandler(carSvc)

	authMw := middleware.NewAuthMiddleware(jwtTestSecret)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.GinBearerJWT(authMw))
	v1 := router.Group("/api/v1")
	v1.GET("/cars/:id", carHandler.GetCar)

	suite.app = router
	suite.db = db
}

func (suite *ClientJWTCarIntegrationSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM cars")
	suite.db.Exec("DELETE FROM users")
}

func (suite *ClientJWTCarIntegrationSuite) TestJWT_ClientCannotGetOtherClientCar() {
	clientA := uuid.New()
	clientB := uuid.New()
	carID := uuid.New()

	suite.db.Create(&postgres.UserModel{
		ID: clientA, Email: "ja@test.com", PasswordHash: "x", FirstName: "A", LastName: "A",
		Role: "client", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	suite.db.Create(&postgres.UserModel{
		ID: clientB, Email: "jb@test.com", PasswordHash: "x", FirstName: "B", LastName: "B",
		Role: "client", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	suite.db.Create(&postgres.CarModel{
		ID: carID, Make: "Ford", Model: "Fiesta", Year: 2018, LicensePlate: "JWT-1", Color: "Blue",
		OwnerID: clientB, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/cars/"+carID.String(), nil)
	req.Header.Set("Authorization", bearerToken(clientA, "client"))
	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func TestClientJWTCarIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ClientJWTCarIntegrationSuite))
}
