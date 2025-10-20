package middleware

import (
	"net/http"
	"strings"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Get auth service from context (you'll need to set this up)
		authService, exists := c.Get("authService")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Auth service not available",
			})
			c.Abort()
			return
		}

		authSvc := authService.(ports.AuthService)
		user, err := authSvc.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user in context for use in handlers
		c.Set("user", user)
		c.Next()
	}
}

func SetAuthService(authService ports.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("authService", authService)
		c.Next()
	}
}
