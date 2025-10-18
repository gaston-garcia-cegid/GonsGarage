package services

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*domain.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	ListEmployees(ctx context.Context, limit, offset int) ([]*domain.Employee, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, req UpdateEmployeeRequest) (*domain.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
}

type CreateEmployeeRequest struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	FirstName   string    `json:"first_name" validate:"required"`
	LastName    string    `json:"last_name" validate:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Position    string    `json:"position" binding:"required"`
	Department  string    `json:"department"`
	PhoneNumber string    `json:"phone_number"`
	HourlyRate  float64   `json:"hourly_rate" validate:"required,min=0"`
}

type UpdateEmployeeRequest struct {
	FirstName   *string  `json:"first_name,omitempty"`
	LastName    *string  `json:"last_name,omitempty"`
	Email       *string  `json:"email"`
	Position    *string  `json:"position"`
	Department  *string  `json:"department"`
	PhoneNumber *string  `json:"phone_number"`
	HourlyRate  *float64 `json:"hourly_rate,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}
