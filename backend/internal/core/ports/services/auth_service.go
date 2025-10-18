package services

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	Register(ctx context.Context, email, password, role string) (*domain.User, error)
	ValidateToken(tokenString string) (*domain.User, error)
	GenerateToken(user *domain.User) (string, error)
}
