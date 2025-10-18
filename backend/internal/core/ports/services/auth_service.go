package services

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*domain.User, string, error)
	Register(ctx context.Context, email, password, role string) (*domain.User, error)
	ValidateToken(tokenString string) (*domain.User, error)
	GenerateToken(user *domain.User) (string, error)
}

type EmployeeService interface {
	CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*domain.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	ListEmployees(ctx context.Context, limit, offset int) ([]*domain.Employee, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, req UpdateEmployeeRequest) (*domain.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
}

type NotificationService interface {
	SendSMS(ctx context.Context, to, message string) error
	SendWhatsApp(ctx context.Context, to, message string) error
	QueueNotification(ctx context.Context, notification NotificationRequest) error
}

// Request/Response DTOs
type CreateEmployeeRequest struct {
	UserID     uuid.UUID `json:"user_id" validate:"required"`
	FirstName  string    `json:"first_name" validate:"required"`
	LastName   string    `json:"last_name" validate:"required"`
	Position   string    `json:"position" validate:"required"`
	HourlyRate float64   `json:"hourly_rate" validate:"required,min=0"`
}

type UpdateEmployeeRequest struct {
	FirstName  *string  `json:"first_name,omitempty"`
	LastName   *string  `json:"last_name,omitempty"`
	Position   *string  `json:"position,omitempty"`
	HourlyRate *float64 `json:"hourly_rate,omitempty"`
	IsActive   *bool    `json:"is_active,omitempty"`
}

type NotificationRequest struct {
	Type    string `json:"type"` // "sms" or "whatsapp"
	To      string `json:"to"`
	Message string `json:"message"`
}
