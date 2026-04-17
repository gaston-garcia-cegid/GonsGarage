package middleware_test

import (
	"reflect"
	"testing"

	legacy "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/middleware"
	appmw "github.com/gaston-garcia-cegid/gonsgarage/internal/middleware"
	"github.com/stretchr/testify/require"
)

func TestMiddleware_functionsPointToLegacy(t *testing.T) {
	cases := []struct {
		name string
		fac  any
		leg  any
	}{
		{"NewAuthMiddleware", appmw.NewAuthMiddleware, legacy.NewAuthMiddleware},
		{"SlogRequestLogger", appmw.SlogRequestLogger, legacy.SlogRequestLogger},
		{"RequireStaffManagers", appmw.RequireStaffManagers, legacy.RequireStaffManagers},
		{"NewIPRateLimiter", appmw.NewIPRateLimiter, legacy.NewIPRateLimiter},
		{"IPRateLimiterFromEnv", appmw.IPRateLimiterFromEnv, legacy.IPRateLimiterFromEnv},
		{"RateLimitAuth", appmw.RateLimitAuth, legacy.RateLimitAuth},
		{"CORSMiddleware", appmw.CORSMiddleware, legacy.CORSMiddleware},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t,
				reflect.ValueOf(tc.leg).Pointer(),
				reflect.ValueOf(tc.fac).Pointer(),
			)
		})
	}
}
