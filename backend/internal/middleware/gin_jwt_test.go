package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGinBearerJWT_ValidTokenSetsContext(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	secret := "unit-test-secret"
	am := NewAuthMiddleware(secret)
	uid := uuid.New()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": uid.String(),
		"role":   "client",
		"email":  "c@example.com",
	})
	signed, err := tok.SignedString([]byte(secret))
	require.NoError(t, err)

	r := gin.New()
	r.Use(GinBearerJWT(am))
	r.GET("/p", func(c *gin.Context) {
		v, _ := c.Get("userID")
		role, _ := c.Get("userRole")
		c.JSON(http.StatusOK, gin.H{"userID": v, "role": role})
	})

	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), uid.String())
	assert.Contains(t, w.Body.String(), "client")
}

func TestGinBearerJWT_MissingHeader(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(GinBearerJWT(NewAuthMiddleware("s")))
	r.GET("/p", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGinBearerJWT_WrongSecret(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userID": uuid.New().String()})
	signed, err := tok.SignedString([]byte("other-secret"))
	require.NoError(t, err)

	r := gin.New()
	r.Use(GinBearerJWT(NewAuthMiddleware("expected-secret")))
	r.GET("/p", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/p", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
