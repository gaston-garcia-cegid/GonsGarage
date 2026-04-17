package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SlogRequestLogger registra cada petición HTTP con slog (JSON o texto según el default global).
// Reduce ruido en /health, /ready y /metrics (solo nivel debug).
func SlogRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		attrs := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency_ms", latency.Milliseconds(),
		}
		if rid := c.GetHeader("X-Request-Id"); rid != "" {
			attrs = append(attrs, "request_id", rid)
		}

		switch {
		case path == "/metrics" || path == "/health" || path == "/ready":
			slog.DebugContext(c.Request.Context(), "request", attrs...)
		case len(c.Errors) > 0:
			attrs = append(attrs, "errors", c.Errors.String())
			slog.ErrorContext(c.Request.Context(), "request", attrs...)
		case status >= 500:
			slog.ErrorContext(c.Request.Context(), "request", attrs...)
		case status >= 400:
			slog.WarnContext(c.Request.Context(), "request", attrs...)
		default:
			slog.InfoContext(c.Request.Context(), "request", attrs...)
		}
	}
}
