package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/gaston-garcia-cegid/gons-garage/backend/internal/core/domain"
	"github.com/gaston-garcia-cegid/gons-garage/backend/internal/core/ports/repositories"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	mock.Mock
}

// NewMockUserRepository creates a new mock user repository
func NewMockUserRepository() repositories.UserRepository {
	return &MockUserRepository{}
}

// Create implements UserRepository.Create
func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetByID implements UserRepository.GetByID
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// GetByEmail implements UserRepository.GetByEmail
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// Update implements UserRepository.Update
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete implements UserRepository.Delete
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List implements UserRepository.List
func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.User), args.Get(1).(int64), args.Error(2)
}

// EmailExists implements UserRepository.EmailExists
func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}
