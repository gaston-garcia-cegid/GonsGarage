package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/middleware"
	authsvc "github.com/gaston-garcia-cegid/gonsgarage/internal/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// provisionTestUserRepo is a minimal UserRepository for admin provision HTTP tests.
type provisionTestUserRepo struct {
	byEmail   map[string]*domain.User
	createErr error
}

func newProvisionTestUserRepo() *provisionTestUserRepo {
	return &provisionTestUserRepo{byEmail: make(map[string]*domain.User)}
}

func (s *provisionTestUserRepo) Create(ctx context.Context, user *domain.User) error {
	if s.createErr != nil {
		return s.createErr
	}
	s.byEmail[user.Email] = user
	return nil
}

func (s *provisionTestUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := s.byEmail[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (s *provisionTestUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	for _, u := range s.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (s *provisionTestUserRepo) GetByRole(ctx context.Context, role string, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (s *provisionTestUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (s *provisionTestUserRepo) Update(ctx context.Context, user *domain.User) error { return nil }
func (s *provisionTestUserRepo) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (s *provisionTestUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return nil
}
func (s *provisionTestUserRepo) GetActiveUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}

func newProvisionTestRouter(t *testing.T, secret string, repo ports.UserRepository) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	am := middleware.NewAuthMiddleware(secret)
	authService := authsvc.NewAuthService(repo, secret, 24)
	h := NewAdminUserHandler(authService)

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(middleware.GinBearerJWT(am))
	admin := api.Group("/admin")
	admin.Use(middleware.RequireStaffManagers())
	admin.POST("/users", h.ProvisionUser)
	return r
}

func TestProvisionUser_NoJWT_Not2xx(t *testing.T) {
	t.Parallel()
	repo := newProvisionTestUserRepo()
	r := newProvisionTestRouter(t, "prov-secret-nojwt", repo)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader([]byte(`{}`)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusCreated, w.Code)
	assert.True(t, w.Code == http.StatusUnauthorized || w.Code == http.StatusBadRequest)
}

func TestProvisionUser_ClientForbidden(t *testing.T) {
	t.Parallel()
	secret := "prov-secret-client"
	repo := newProvisionTestUserRepo()
	r := newProvisionTestRouter(t, secret, repo)
	uid := uuid.New()

	body := map[string]string{
		"email": "new@example.com", "password": "secret12", "firstName": "A", "lastName": "B", "role": domain.RoleClient,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleClient))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProvisionUser_EmployeeForbidden(t *testing.T) {
	t.Parallel()
	secret := "prov-secret-emp"
	repo := newProvisionTestUserRepo()
	r := newProvisionTestRouter(t, secret, repo)
	uid := uuid.New()

	body := map[string]string{
		"email": "new2@example.com", "password": "secret12", "firstName": "A", "lastName": "B", "role": domain.RoleClient,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleEmployee))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProvisionUser_AdminBodyAdminRole_Not2xx(t *testing.T) {
	t.Parallel()
	secret := "prov-secret-adm-admin"
	repo := newProvisionTestUserRepo()
	r := newProvisionTestRouter(t, secret, repo)
	uid := uuid.New()

	body := map[string]string{
		"email": "adm@example.com", "password": "secret12", "firstName": "A", "lastName": "B", "role": domain.RoleAdmin,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleAdmin))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProvisionUser_AdminCreatesClient_201(t *testing.T) {
	t.Parallel()
	secret := "prov-secret-adm-ok"
	repo := newProvisionTestUserRepo()
	r := newProvisionTestRouter(t, secret, repo)
	uid := uuid.New()

	body := map[string]string{
		"email": "clientnew@example.com", "password": "secret12", "firstName": "A", "lastName": "B", "role": domain.RoleClient,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleAdmin))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var out struct {
		User domain.User `json:"user"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	assert.Equal(t, domain.RoleClient, out.User.Role)
	assert.Equal(t, "clientnew@example.com", out.User.Email)
	assert.Empty(t, out.User.Password)
}

func TestProvisionUser_ManagerCreatesManager_Not2xx(t *testing.T) {
	t.Parallel()
	secret := "prov-secret-mgr-mgr"
	repo := newProvisionTestUserRepo()
	r := newProvisionTestRouter(t, secret, repo)
	uid := uuid.New()

	body := map[string]string{
		"email": "mgrdup@example.com", "password": "secret12", "firstName": "A", "lastName": "B", "role": domain.RoleManager,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleManager))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestProvisionUser_ManagerCreatesEmployee_201(t *testing.T) {
	t.Parallel()
	secret := "prov-secret-mgr-emp"
	repo := newProvisionTestUserRepo()
	r := newProvisionTestRouter(t, secret, repo)
	uid := uuid.New()

	body := map[string]string{
		"email": "empnew@example.com", "password": "secret12", "firstName": "A", "lastName": "B", "role": domain.RoleEmployee,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+testJWT(t, secret, uid, domain.RoleManager))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var out struct {
		User domain.User `json:"user"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	assert.Equal(t, domain.RoleEmployee, out.User.Role)
}
