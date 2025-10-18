package repositories

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
)

// EmployeeRepository defines employee data operations
type EmployeeRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, employee *domain.Employee) error
	GetByID(ctx context.Context, id string) (*domain.Employee, error)
	Update(ctx context.Context, employee *domain.Employee) error
	Delete(ctx context.Context, id string) error

	// Business-specific queries
	GetByDepartment(ctx context.Context, department string) ([]*domain.Employee, error)
	GetActiveEmployees(ctx context.Context) ([]*domain.Employee, error)
	GetEmployeesByRole(ctx context.Context, role string) ([]*domain.Employee, error)

	// Pagination and filtering
	List(ctx context.Context, filters *domain.EmployeeFilters) ([]*domain.Employee, int64, error)

	// Reporting queries
	GetEmployeeWorkHours(ctx context.Context, employeeID string, from, to time.Time) ([]*domain.WorkHour, error)
	GetEmployeesOnVacation(ctx context.Context, date time.Time) ([]*domain.Employee, error)

	// Search functionality
	SearchByName(ctx context.Context, name string, limit int) ([]*domain.Employee, error)
}

// EmployeeFilters for query filtering
type EmployeeFilters struct {
	Department *string
	Role       *string
	IsActive   *bool
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  string
}
