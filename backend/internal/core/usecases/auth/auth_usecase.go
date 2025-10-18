package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo   ports.UserRepository
	jwtSecret  string
	expireTime time.Duration
}

func NewAuthUseCase(userRepo ports.UserRepository, jwtSecret string, expireHours int) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		jwtSecret:  jwtSecret,
		expireTime: time.Duration(expireHours) * time.Hour,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.ValidatePassword(password) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := uc.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (uc *AuthUseCase) Register(ctx context.Context, email, password, role string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	user, err := domain.NewUser(email, password, role)
	if err != nil {
		return nil, err
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) GenerateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(uc.expireTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user, err := uc.userRepo.GetByID(context.Background(), userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
