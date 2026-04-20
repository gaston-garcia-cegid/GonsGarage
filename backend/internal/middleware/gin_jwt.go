package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// jwtClaimString coerces JWT MapClaims values that may decode as string, json.Number, or float64.
func jwtClaimString(v interface{}) (string, bool) {
	switch t := v.(type) {
	case string:
		return t, true
	case json.Number:
		return t.String(), true
	default:
		return "", false
	}
}

// GinBearerJWT validates Authorization: Bearer <JWT> and sets userID (string), userRole, userEmail on Gin context.
// Mirrors production auth used by API handlers (see cmd/api).
func GinBearerJWT(auth *AuthMiddleware) gin.HandlerFunc {
	secret := auth.GetJWTSecret()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		var userIDStr string
		if uid, exists := claims["userID"]; exists {
			if uidStr, ok := uid.(string); ok {
				userIDStr = uidStr
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID in token"})
				c.Abort()
				return
			}
		} else if sub, exists := claims["sub"]; exists {
			if subStr, ok := sub.(string); ok {
				userIDStr = subStr
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid sub in token"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user identifier in token"})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID format"})
			c.Abort()
			return
		}

		c.Set("userID", userID.String())

		if email, exists := claims["email"]; exists {
			if emailStr, ok := jwtClaimString(email); ok && emailStr != "" {
				c.Set("userEmail", emailStr)
			}
		}
		if role, exists := claims["role"]; exists {
			if roleStr, ok := jwtClaimString(role); ok && roleStr != "" {
				c.Set("userRole", roleStr)
			}
		}

		log.Printf("✅ Authentication successful for user: %s", userID.String())
		c.Next()
	}
}
