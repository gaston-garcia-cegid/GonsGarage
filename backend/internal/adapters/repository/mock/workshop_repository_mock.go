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

// CreateWorkshop implements WorkshopRepository.CreateWorkshop
func (m *MockWorkshopRepository) CreateWorkshop(ctx context.Context, workshop *domain.Workshop) error {
	args := m.Called(ctx, workshop)
	return args.Error(0)
}

// GetWorkshopByID implements WorkshopRepository.GetWorkshopByID
func (m *MockWorkshopRepository) GetWorkshopByID(ctx context.Context, id string) (*domain.Workshop, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Workshop), args.Error(1)
}

// UpdateWorkshop implements WorkshopRepository.UpdateWorkshop
func (m *MockWorkshopRepository) UpdateWorkshop(ctx context.Context, workshop *domain.Workshop) error {
	args := m.Called(ctx, workshop)
	return args.Error(0)
}

// DeleteWorkshop implements WorkshopRepository.DeleteWorkshop
func (m *MockWorkshopRepository) DeleteWorkshop(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// ListWorkshops implements WorkshopRepository.ListWorkshops
func (m *MockWorkshopRepository) ListWorkshops(ctx context.Context, limit, offset int) ([]*domain.Workshop, int64, error) {
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
