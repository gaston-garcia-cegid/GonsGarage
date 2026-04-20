package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testJWTWorkshop(t *testing.T, secret string, userID uuid.UUID, role string) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID.String(),
		"role":   role,
	})
	s, err := tok.SignedString([]byte(secret))
	require.NoError(t, err)
	return s
}

func TestRequireWorkshopStaff_ClientForbidden(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "ws-staff-secret"
	am := NewAuthMiddleware(secret)
	uid := uuid.New()

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(GinBearerJWT(am))
	api.Use(RequireWorkshopStaff())
	api.GET("/suppliers", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers", nil)
	req.Header.Set("Authorization", "Bearer "+testJWTWorkshop(t, secret, uid, domain.RoleClient))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireWorkshopStaff_EmployeeOK(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "ws-staff-secret-emp"
	am := NewAuthMiddleware(secret)
	uid := uuid.New()

	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(GinBearerJWT(am))
	api.Use(RequireWorkshopStaff())
	api.GET("/suppliers", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/suppliers", nil)
	req.Header.Set("Authorization", "Bearer "+testJWTWorkshop(t, secret, uid, domain.RoleEmployee))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
