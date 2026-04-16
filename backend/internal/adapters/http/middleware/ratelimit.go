package middleware

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// NewIPRateLimiter returns a simple per-IP token bucket (golang.org/x/time/rate).
func NewIPRateLimiter(rps float64, burst int) *IPRateLimiter {
	if rps <= 0 {
		rps = 5
	}
	if burst <= 0 {
		burst = 10
	}
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        rate.Limit(rps),
		b:        burst,
	}
}

// IPRateLimiterFromEnv builds a limiter from AUTH_RATE_LIMIT_RPS and AUTH_RATE_LIMIT_BURST (defaults 5 / 10).
func IPRateLimiterFromEnv() *IPRateLimiter {
	rps := 5.0
	if v := strings.TrimSpace(os.Getenv("AUTH_RATE_LIMIT_RPS")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			rps = f
		}
	}
	burst := 10
	if v := strings.TrimSpace(os.Getenv("AUTH_RATE_LIMIT_BURST")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			burst = n
		}
	}
	return NewIPRateLimiter(rps, burst)
}

// IPRateLimiter tracks one limiter per client IP.
type IPRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	b        int
}

func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	lim, ok := i.limiters[ip]
	if !ok {
		lim = rate.NewLimiter(i.r, i.b)
		i.limiters[ip] = lim
	}
	return lim
}

// RateLimitAuth returns middleware that returns 429 when the bucket for ClientIP is empty.
func RateLimitAuth(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}
