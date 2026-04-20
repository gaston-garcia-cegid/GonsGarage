package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/middleware"
	repairsvc "github.com/gaston-garcia-cegid/gonsgarage/internal/service/repair"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testJWT(t *testing.T, secret string, userID uuid.UUID, role string) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID.String(),
		"role":   role,
		"email":  role + "@test.local",
	})
	s, err := tok.SignedString([]byte(secret))
	require.NoError(t, err)
	return s
}

// --- Employees: RequireStaffManagers (spec client + employee forbidden; manager allowed) ---

func TestMVPAccess_EmployeesGET_ClientForbidden(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "mvp-role-access-secret"
	am := middleware.NewAuthMiddleware(secret)
	uid := uuid.New()

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(middleware.GinBearerJWT(am))
	employees := api.Group("/employees")
	employees.Use(middleware.RequireStaffManagers())
	employees.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleClient))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMVPAccess_EmployeesGET_EmployeeForbidden(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "mvp-role-access-secret-emp"
	am := middleware.NewAuthMiddleware(secret)
	uid := uuid.New()

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(middleware.GinBearerJWT(am))
	employees := api.Group("/employees")
	employees.Use(middleware.RequireStaffManagers())
	employees.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleEmployee))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMVPAccess_EmployeesGET_ManagerReachesHandler(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "mvp-role-access-secret-mgr"
	am := middleware.NewAuthMiddleware(secret)
	uid := uuid.New()

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(middleware.GinBearerJWT(am))
	employees := api.Group("/employees")
	employees.Use(middleware.RequireStaffManagers())
	employees.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleManager))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestMVPAccess_EmployeesGET_AdminReachesHandler(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "mvp-role-access-secret-adm"
	am := middleware.NewAuthMiddleware(secret)
	uid := uuid.New()

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(middleware.GinBearerJWT(am))
	employees := api.Group("/employees")
	employees.Use(middleware.RequireStaffManagers())
	employees.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/employees", nil)
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleAdmin))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

// --- Repairs POST: client 403; employee 201 with stubs ---

type mvpUserRepo struct {
	byID map[uuid.UUID]*domain.User
}

func (m *mvpUserRepo) Create(context.Context, *domain.User) error                     { return errors.New("mvpUserRepo.Create not used") }
func (m *mvpUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return u, nil
}
func (m *mvpUserRepo) GetByEmail(context.Context, string) (*domain.User, error) { return nil, domain.ErrUserNotFound }
func (m *mvpUserRepo) GetByRole(context.Context, string, int, int) ([]*domain.User, error) {
	return nil, nil
}
func (m *mvpUserRepo) List(context.Context, int, int) ([]*domain.User, error) { return nil, nil }
func (m *mvpUserRepo) Update(context.Context, *domain.User) error             { return errors.New("not used") }
func (m *mvpUserRepo) Delete(context.Context, uuid.UUID) error                  { return errors.New("not used") }
func (m *mvpUserRepo) UpdatePassword(context.Context, uuid.UUID, string) error {
	return errors.New("not used")
}
func (m *mvpUserRepo) GetActiveUsers(context.Context, int, int) ([]*domain.User, error) {
	return nil, nil
}

type mvpCarRepo struct {
	car *domain.Car
}

func (m *mvpCarRepo) Create(context.Context, *domain.Car) error { return errors.New("not used") }
func (m *mvpCarRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.Car, error) {
	if m.car != nil && m.car.ID == id {
		return m.car, nil
	}
	return nil, fmt.Errorf("car not found")
}
func (m *mvpCarRepo) GetByOwnerID(context.Context, uuid.UUID) ([]*domain.Car, error) {
	return nil, nil
}
func (m *mvpCarRepo) GetByLicensePlate(context.Context, string) (*domain.Car, error) {
	return nil, domain.ErrUserNotFound
}
func (m *mvpCarRepo) List(context.Context, int, int) ([]*domain.Car, error) { return nil, nil }
func (m *mvpCarRepo) Update(context.Context, *domain.Car) error             { return errors.New("not used") }
func (m *mvpCarRepo) Delete(context.Context, uuid.UUID) error                { return errors.New("not used") }
func (m *mvpCarRepo) GetWithRepairs(context.Context, uuid.UUID) (*domain.Car, error) {
	return nil, nil
}
func (m *mvpCarRepo) GetDeletedByLicensePlate(context.Context, string) (*domain.Car, error) {
	return nil, nil
}
func (m *mvpCarRepo) Restore(context.Context, uuid.UUID) error { return errors.New("not used") }

type mvpRepairRepo struct {
	byCar map[uuid.UUID][]*domain.Repair
}

func (m *mvpRepairRepo) Create(context.Context, *domain.Repair) error { return nil }
func (m *mvpRepairRepo) GetByID(context.Context, uuid.UUID) (*domain.Repair, error) {
	return nil, domain.ErrRepairNotFound
}
func (m *mvpRepairRepo) Update(context.Context, *domain.Repair) error { return errors.New("not used") }
func (m *mvpRepairRepo) Delete(context.Context, uuid.UUID) error      { return errors.New("not used") }
func (m *mvpRepairRepo) GetByCarID(_ context.Context, carID uuid.UUID) ([]*domain.Repair, error) {
	if m.byCar == nil {
		return nil, nil
	}
	return m.byCar[carID], nil
}

var _ ports.UserRepository = (*mvpUserRepo)(nil)
var _ ports.CarRepository = (*mvpCarRepo)(nil)
var _ ports.RepairRepository = (*mvpRepairRepo)(nil)

func repairPOSTRouter(t *testing.T, secret string, userRepo ports.UserRepository, carRepo ports.CarRepository, repRepo ports.RepairRepository) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	am := middleware.NewAuthMiddleware(secret)
	svc := repairsvc.NewRepairService(repRepo, carRepo, userRepo)
	h := NewRepairHandler(svc)

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(middleware.GinBearerJWT(am))
	api.POST("/repairs", h.GinCreateRepair)
	return r
}

func TestMVPAccess_RepairPOST_ClientForbidden(t *testing.T) {
	t.Parallel()
	secret := "mvp-repair-client-secret"
	clientID := uuid.New()
	carID := uuid.New()

	u, err := domain.NewUser("c@test.local", "x", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	u.ID = clientID

	users := &mvpUserRepo{byID: map[uuid.UUID]*domain.User{clientID: u}}
	cars := &mvpCarRepo{car: &domain.Car{ID: carID, OwnerID: clientID}}
	repairs := &mvpRepairRepo{byCar: map[uuid.UUID][]*domain.Repair{carID: {}}}

	r := repairPOSTRouter(t, secret, users, cars, repairs)

	body := map[string]any{
		"car_id":      carID.String(),
		"description": "Brake inspection",
		"status":      "pending",
		"cost":        0.0,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/repairs", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, clientID, domain.RoleClient))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMVPAccess_RepairPOST_EmployeeCreated(t *testing.T) {
	t.Parallel()
	secret := "mvp-repair-emp-secret"
	empID := uuid.New()
	carID := uuid.New()
	ownerID := uuid.New()

	emp, err := domain.NewUser("e@test.local", "x", "E", "M", domain.RoleEmployee)
	require.NoError(t, err)
	emp.ID = empID

	users := &mvpUserRepo{byID: map[uuid.UUID]*domain.User{empID: emp}}
	cars := &mvpCarRepo{car: &domain.Car{ID: carID, OwnerID: ownerID}}
	repairs := &mvpRepairRepo{byCar: map[uuid.UUID][]*domain.Repair{carID: {}}}

	r := repairPOSTRouter(t, secret, users, cars, repairs)

	body := map[string]any{
		"car_id":      carID.String(),
		"description": "Oil change",
		"status":      "pending",
		"cost":        10.0,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/repairs", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, empID, domain.RoleEmployee))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}
