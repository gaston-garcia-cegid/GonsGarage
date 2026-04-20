package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextUserID returns the authenticated user UUID from Gin context (JWT middleware).
func ContextUserID(c *gin.Context) (uuid.UUID, error) {
	v, ok := c.Get("userID")
	if !ok {
		return uuid.Nil, fmt.Errorf("missing userID")
	}
	s, ok := v.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid userID type")
	}
	return uuid.Parse(s)
}

// QueryLimitOffset parses limit and offset query params with safe defaults.
func QueryLimitOffset(c *gin.Context, defaultLimit, maxLimit int) (limit, offset int) {
	limit = defaultLimit
	if q := strings.TrimSpace(c.Query("limit")); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	offset = 0
	if q := strings.TrimSpace(c.Query("offset")); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}
