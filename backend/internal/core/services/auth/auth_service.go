package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo   ports.UserRepository
	jwtSecret  string
	expireTime time.Duration
}

func NewAuthService(userRepo ports.UserRepository, jwtSecret string, expireHours int) ports.AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		expireTime: time.Duration(expireHours) * time.Hour,
	}
}

func (uc *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}

	if !user.ValidatePassword(password) {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", errors.New("user account is deactivated")
	}

	token, _, err := uc.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *AuthService) Register(ctx context.Context, req ports.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Validate role
	if req.Role == "" {
		req.Role = "employee"
	}

	// Validate role is one of the allowed values
	allowedRoles := []string{"admin", "manager", "employee", "client"}
	validRole := false
	for _, role := range allowedRoles {
		if req.Role == role {
			validRole = true
			break
		}
	}
	if !validRole {
		return nil, errors.New("invalid role. Must be one of: admin, manager, employee, client")
	}

	user, err := domain.NewUser(req.Email, req.Password, req.FirstName, req.LastName, req.Role)
	if err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""

	return user, nil
}

func (uc *AuthService) GenerateToken(user *domain.User) (string, time.Time, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key"
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"userID":     user.ID.String(),
		"sub":        user.ID.String(),
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"exp":        time.Now().Add(uc.expireTime).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (uc *AuthService) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["userID"].(string)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil, errors.New("invalid user ID in token")
		}

		user, err := uc.userRepo.GetByID(context.Background(), userID)
		if err != nil {
			return nil, err
		}

		if !user.IsActive {
			return nil, errors.New("user account is deactivated")
		}

		return user, nil
	}

	return nil, errors.New("invalid token")
}

func (uc *AuthService) RefreshToken(ctx context.Context, token string) (string, error) {
	user, err := uc.ValidateToken(token)
	if err != nil {
		return "", err
	}

	token, _, err = uc.GenerateToken(user)
	return token, err
}
