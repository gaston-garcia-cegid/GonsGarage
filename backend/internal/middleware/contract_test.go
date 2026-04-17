package middleware_test

import (
	"testing"

	appmw "github.com/gaston-garcia-cegid/gonsgarage/internal/middleware"
	"github.com/stretchr/testify/require"
)

func TestNewAuthMiddleware_nonNil(t *testing.T) {
	m := appmw.NewAuthMiddleware("test-secret")
	require.NotNil(t, m)
	require.Equal(t, "test-secret", m.GetJWTSecret())
}

func TestIPRateLimiterFromEnv_nonNil(t *testing.T) {
	l := appmw.IPRateLimiterFromEnv()
	require.NotNil(t, l)
}
