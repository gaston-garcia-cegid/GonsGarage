// Package middleware hosts Gin/HTTP middleware (template §1).
// Phase 1: aliases to internal/adapters/http/middleware.
package middleware

import legacy "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/http/middleware"

type (
	AuthMiddleware = legacy.AuthMiddleware
	IPRateLimiter  = legacy.IPRateLimiter
)

var (
	NewAuthMiddleware    = legacy.NewAuthMiddleware
	SlogRequestLogger    = legacy.SlogRequestLogger
	RequireStaffManagers = legacy.RequireStaffManagers
	NewIPRateLimiter     = legacy.NewIPRateLimiter
	IPRateLimiterFromEnv = legacy.IPRateLimiterFromEnv
	RateLimitAuth        = legacy.RateLimitAuth
	CORSMiddleware       = legacy.CORSMiddleware
)
