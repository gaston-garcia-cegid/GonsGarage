package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
	}
}

// ✅ Fixed: Return http.Handler for proper middleware chaining
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if token has Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil {
			fmt.Printf("Token parsing error: %v\n", err) // Debug log
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			fmt.Printf("Token is not valid\n") // Debug log
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// ✅ Fixed: Extract claims as MapClaims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Printf("Cannot parse token claims\n") // Debug log
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		var userIDStr string
		if uid, exists := claims["userID"]; exists {
			if uidStr, ok := uid.(string); ok {
				userIDStr = uidStr
			} else {
				fmt.Printf("userID claim is not a string: %T\n", uid)
				http.Error(w, "Invalid userID in token", http.StatusUnauthorized)
				return
			}
		} else if sub, exists := claims["sub"]; exists {
			if subStr, ok := sub.(string); ok {
				userIDStr = subStr
			} else {
				fmt.Printf("sub claim is not a string: %T\n", sub)
				http.Error(w, "Invalid sub in token", http.StatusUnauthorized)
				return
			}
		} else {
			fmt.Printf("No userID or sub claim found in token\n")
			http.Error(w, "Missing user identifier in token", http.StatusUnauthorized)
			return
		}

		// Validate UUID format
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			fmt.Printf("Invalid UUID format for userID: %s, error: %v\n", userIDStr, err)
			http.Error(w, "Invalid userID format", http.StatusUnauthorized)
			return
		}

		// ✅ Fixed: Store userID as string to match handler expectations
		ctx := context.WithValue(r.Context(), "userID", userID.String())

		// Add user info to request context
		if email, exists := claims["email"]; exists {
			if emailStr, ok := email.(string); ok {
				ctx = context.WithValue(ctx, "userEmail", emailStr)
			}
		}
		if role, exists := claims["role"]; exists {
			if roleStr, ok := role.(string); ok {
				ctx = context.WithValue(ctx, "userRole", roleStr)
			}
		}

		fmt.Printf("✅ Authentication successful for user: %s\n", userID.String()) // Debug log

		// ✅ Continue with the request using updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userIDStr, ok := ctx.Value("userID").(string)
	if !ok {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, false
	}

	return userID, true
}

// GetUserEmailFromContext extracts user email from request context
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value("userEmail").(string)
	return email, ok
}

// GetUserRoleFromContext extracts user role from request context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value("userRole").(string)
	return role, ok
}

func (m *AuthMiddleware) GetJWTSecret() string {
	return m.jwtSecret
}
