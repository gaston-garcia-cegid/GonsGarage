package repositories

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EmployeeRepository interface {
	Create(ctx context.Context, employee *domain.Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Employee, error)
	Update(ctx context.Context, employee *domain.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
