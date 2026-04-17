//go:build cgo

package integration

import (
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

	"github.com/gaston-garcia-cegid/gonsgarage/internal/handler"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/repository/postgres"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/service/appointment"
)

type AppointmentIntegrationTestSuite struct {
	suite.Suite
	app *gin.Engine
	db  *gorm.DB
}

func (suite *AppointmentIntegrationTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: nil})
	require.NoError(suite.T(), err)

	err = db.AutoMigrate(&postgres.UserModel{}, &postgres.CarModel{}, &postgres.AppointmentModel{})
	require.NoError(suite.T(), err)

	userRepo := postgres.NewPostgresUserRepository(db)
	carRepo := postgres.NewPostgresCarRepository(db)
	apptRepo := postgres.NewPostgresAppointmentRepository(db)
	apptSvc := appointment.NewAppointmentService(apptRepo, userRepo, carRepo)
	apptHandler := handler.NewAppointmentHandler(apptSvc)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(suite.testAuthMiddleware())
	v1 := router.Group("/api/v1")
	v1.GET("/appointments/:id", apptHandler.GetAppointment)

	suite.app = router
	suite.db = db
}

func (suite *AppointmentIntegrationTestSuite) TearDownTest() {
	suite.db.Exec("DELETE FROM appointments")
	suite.db.Exec("DELETE FROM cars")
	suite.db.Exec("DELETE FROM users")
}

func (suite *AppointmentIntegrationTestSuite) testAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if uid := c.GetHeader("X-User-ID"); uid != "" {
			if parsed, err := uuid.Parse(uid); err == nil {
				c.Set("userID", parsed.String())
			}
		}
		c.Next()
	}
}

func (suite *AppointmentIntegrationTestSuite) TestGetAppointment_ClientForbiddenOtherCustomer() {
	clientA := uuid.New()
	clientB := uuid.New()
	carID := uuid.New()
	apptID := uuid.New()

	suite.db.Create(&postgres.UserModel{
		ID: clientA, Email: "a@test.com", PasswordHash: "x", FirstName: "A", LastName: "A",
		Role: "client", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	suite.db.Create(&postgres.UserModel{
		ID: clientB, Email: "b@test.com", PasswordHash: "x", FirstName: "B", LastName: "B",
		Role: "client", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	suite.db.Create(&postgres.CarModel{
		ID: carID, Make: "VW", Model: "Golf", Year: 2020, LicensePlate: "AP-1", Color: "Red",
		OwnerID: clientB, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	suite.db.Create(&postgres.AppointmentModel{
		ID: apptID, CustomerID: clientB, CarID: carID, ScheduledTime: time.Now().Add(time.Hour),
		Notes: "", Status: "scheduled", ServiceType: "service", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/appointments/"+apptID.String(), nil)
	req.Header.Set("X-User-ID", clientA.String())
	w := httptest.NewRecorder()
	suite.app.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	var body map[string]interface{}
	require.NoError(suite.T(), json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(suite.T(), "forbidden", body["error"])
}

func TestAppointmentIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AppointmentIntegrationTestSuite))
}
