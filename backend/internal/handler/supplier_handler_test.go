package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubSupplierService struct {
	list []*domain.Supplier
	tot  int64
}

func (s *stubSupplierService) Create(ctx context.Context, row *domain.Supplier, requestingUserID uuid.UUID) (*domain.Supplier, error) {
	id := uuid.New()
	return &domain.Supplier{
		ID: id, Name: row.Name, IsActive: row.IsActive, ContactEmail: row.ContactEmail,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}, nil
}

func (s *stubSupplierService) Get(ctx context.Context, id uuid.UUID, requestingUserID uuid.UUID) (*domain.Supplier, error) {
	return nil, domain.ErrSupplierNotFound
}

func (s *stubSupplierService) List(ctx context.Context, requestingUserID uuid.UUID, limit, offset int) ([]*domain.Supplier, int64, error) {
	return s.list, s.tot, nil
}

func (s *stubSupplierService) Update(ctx context.Context, row *domain.Supplier, requestingUserID uuid.UUID) (*domain.Supplier, error) {
	return nil, domain.ErrSupplierNotFound
}

func (s *stubSupplierService) Delete(ctx context.Context, id uuid.UUID, requestingUserID uuid.UUID) error {
	return nil
}

func testJWTHandler(t *testing.T, secret string, userID uuid.UUID, role string) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID.String(),
		"role":   role,
	})
	out, err := tok.SignedString([]byte(secret))
	require.NoError(t, err)
	return out
}

func TestSupplierHandler_ListSuppliers_Integration(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "supplier-handler-secret"
	am := middleware.NewAuthMiddleware(secret)
	uid := uuid.New()

	h := NewSupplierHandler(&stubSupplierService{
		list: []*domain.Supplier{{ID: uuid.New(), Name: "ACME", IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}},
		tot:  1,
	})
	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(middleware.GinBearerJWT(am))
	sg := api.Group("/suppliers")
	sg.Use(middleware.RequireWorkshopStaff())
	sg.GET("", h.ListSuppliers)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+testJWTHandler(t, secret, uid, domain.RoleManager))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.EqualValues(t, 1, body["total"])
	items, ok := body["items"].([]interface{})
	require.True(t, ok)
	require.Len(t, items, 1)
}
