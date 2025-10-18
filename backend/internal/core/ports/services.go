package ports

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

// AuthService define os métodos do serviço de autenticação
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, req RegisterRequest) (*domain.User, error)
	ValidateToken(token string) (*domain.User, error)
	RefreshToken(ctx context.Context, token string) (string, error)
}

// RegisterRequest representa os dados para registro
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

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
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name" gorm:"not null"`
	Email        string    `json:"email" binding:"required,email"`
	Position     string    `json:"position" binding:"required"`
	HourlyRate   float64   `json:"hourly_rate" gorm:"not null"`
	HoursWorked  float64   `json:"hours_worked" gorm:"default:0"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	EmployeeCode string    `json:"employee_code"`
	Department   string    `json:"department"`
	HireDate     time.Time `json:"hire_date"`
	Salary       float64   `json:"salary"`
	HoursPerWeek int       `json:"hours_per_week"`
	PhoneNumber  string    `json:"phone_number"`
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

// NotificationService defines the interface for the notification service
type NotificationService interface {
	SendWorkshopCreatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopUpdatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopDeletedNotification(ctx context.Context, workshop *domain.Workshop) error
}

// AccountingService defines the interface for the accounting service
type AccountingService interface {
	CreateAccountingEntry(ctx context.Context, entry *domain.AccountingEntry) error
	UpdateAccountingEntry(ctx context.Context, entry *domain.AccountingEntry) error
	DeleteAccountingEntry(ctx context.Context, id string) error
	ListAccountingEntries(ctx context.Context, limit, offset int) ([]*domain.AccountingEntry, int64, error)
}
