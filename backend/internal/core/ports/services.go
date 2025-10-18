package ports

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

// EmployeeService define os métodos do serviço de funcionários
type EmployeeService interface {
	CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*domain.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	ListEmployees(ctx context.Context, limit, offset int) ([]*domain.Employee, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, req UpdateEmployeeRequest) (*domain.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
}

// CreateEmployeeRequest representa os dados para criar um funcionário
type CreateEmployeeRequest struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Position    string `json:"position" binding:"required"`
	Department  string `json:"department"`
	PhoneNumber string `json:"phone_number"`
}

// UpdateEmployeeRequest representa os dados para atualizar um funcionário
type UpdateEmployeeRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	Position    *string `json:"position"`
	Department  *string `json:"department"`
	PhoneNumber *string `json:"phone_number"`
}

// WorkshopService defines the interface for the workshop service
type WorkshopService interface {
	CreateWorkshop(ctx context.Context, workshop *domain.Workshop) error
	UpdateWorkshop(ctx context.Context, workshop *domain.Workshop) error
	DeleteWorkshop(ctx context.Context, id string) error
	GetWorkshopByID(ctx context.Context, id string) (*domain.Workshop, error)
	ListWorkshops(ctx context.Context, limit, offset int) ([]*domain.Workshop, int64, error)
}
