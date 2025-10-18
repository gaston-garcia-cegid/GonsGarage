package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports/repositories"
)

// MockWorkshopRepository is a mock implementation of WorkshopRepository for testing
type MockWorkshopRepository struct {
	mock.Mock
}

// NewMockWorkshopRepository creates a new mock workshop repository
func NewMockWorkshopRepository() repositories.WorkshopRepository {
	return &MockWorkshopRepository{}
}

// Create implements WorkshopRepository.Create
func (m *MockWorkshopRepository) Create(ctx context.Context, workshop *domain.Workshop) error {
	args := m.Called(ctx, workshop)
	return args.Error(0)
}

// GetByID implements WorkshopRepository.GetByID
func (m *MockWorkshopRepository) GetByID(ctx context.Context, id string) (*domain.Workshop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Workshop), args.Error(1)
}

// Update implements WorkshopRepository.Update
func (m *MockWorkshopRepository) Update(ctx context.Context, workshop *domain.Workshop) error {
	args := m.Called(ctx, workshop)
	return args.Error(0)
}

// Delete implements WorkshopRepository.Delete
func (m *MockWorkshopRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List implements WorkshopRepository.List
func (m *MockWorkshopRepository) List(ctx context.Context, limit, offset int) ([]*domain.Workshop, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Workshop), args.Get(1).(int64), args.Error(2)
}

// SearchByName implements WorkshopRepository.SearchByName
func (m *MockWorkshopRepository) SearchByName(ctx context.Context, name string, limit int) ([]*domain.Workshop, error) {
	args := m.Called(ctx, name, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Workshop), args.Error(1)
}
