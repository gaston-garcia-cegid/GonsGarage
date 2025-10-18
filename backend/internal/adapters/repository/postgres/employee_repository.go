package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/core/ports/repositories"
)

// PostgresEmployeeRepository implements EmployeeRepository interface
type PostgresEmployeeRepository struct {
	db *gorm.DB
}

// NewPostgresEmployeeRepository creates a new PostgreSQL employee repository
func NewPostgresEmployeeRepository(db *gorm.DB) repositories.EmployeeRepository {
	return &PostgresEmployeeRepository{
		db: db,
	}
}

// Create implements EmployeeRepository.Create
func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	dbEmployee := &EmployeeModel{
		ID:           employee.ID,
		UserID:       employee.UserID,
		EmployeeCode: employee.EmployeeCode,
		Position:     employee.Position,
		Department:   employee.Department,
		HireDate:     employee.HireDate,
		Salary:       employee.Salary,
		HoursPerWeek: employee.HoursPerWeek,
		IsActive:     employee.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(dbEmployee).Error; err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}

	employee.CreatedAt = dbEmployee.CreatedAt
	employee.UpdatedAt = dbEmployee.UpdatedAt

	return nil
}

// GetByID implements EmployeeRepository.GetByID
func (r *PostgresEmployeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	var dbEmployee EmployeeModel

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("employees.id = ? AND employees.deleted_at IS NULL", id).
		First(&dbEmployee).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee by ID: %w", err)
	}

	return r.toDomainEmployee(&dbEmployee), nil
}

// GetByDepartment implements EmployeeRepository.GetByDepartment
func (r *PostgresEmployeeRepository) GetByDepartment(ctx context.Context, department string) ([]*domain.Employee, error) {
	var dbEmployees []EmployeeModel

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("department = ? AND deleted_at IS NULL", department).
		Find(&dbEmployees).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get employees by department: %w", err)
	}

	employees := make([]*domain.Employee, len(dbEmployees))
	for i, dbEmployee := range dbEmployees {
		employees[i] = r.toDomainEmployee(&dbEmployee)
	}

	return employees, nil
}

// GetActiveEmployees implements EmployeeRepository.GetActiveEmployees
func (r *PostgresEmployeeRepository) GetActiveEmployees(ctx context.Context) ([]*domain.Employee, error) {
	var dbEmployees []EmployeeModel

	err := r.db.WithContext(ctx).
		Preload("User").
		Where("is_active = ? AND deleted_at IS NULL", true).
		Find(&dbEmployees).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active employees: %w", err)
	}

	employees := make([]*domain.Employee, len(dbEmployees))
	for i, dbEmployee := range dbEmployees {
		employees[i] = r.toDomainEmployee(&dbEmployee)
	}

	return employees, nil
}

// Update implements EmployeeRepository.Update
func (r *PostgresEmployeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	dbEmployee := &EmployeeModel{
		ID:           employee.ID,
		UserID:       employee.UserID,
		EmployeeCode: employee.EmployeeCode,
		Position:     employee.Position,
		Department:   employee.Department,
		HireDate:     employee.HireDate,
		Salary:       employee.Salary,
		HoursPerWeek: employee.HoursPerWeek,
		IsActive:     employee.IsActive,
		UpdatedAt:    time.Now(),
	}

	result := r.db.WithContext(ctx).Model(dbEmployee).Where("id = ?", employee.ID).Updates(dbEmployee)
	if result.Error != nil {
		return fmt.Errorf("failed to update employee: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrEmployeeNotFound
	}

	employee.UpdatedAt = dbEmployee.UpdatedAt
	return nil
}

// Delete implements EmployeeRepository.Delete
func (r *PostgresEmployeeRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Model(&EmployeeModel{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrEmployeeNotFound
	}

	return nil
}

// List implements EmployeeRepository.List
func (r *PostgresEmployeeRepository) List(ctx context.Context, filters *repositories.EmployeeFilters) ([]*domain.Employee, int64, error) {
	query := r.db.WithContext(ctx).Model(&EmployeeModel{}).
		Preload("User").
		Where("employees.deleted_at IS NULL")

	// Apply filters
	if filters.Department != nil {
		query = query.Where("department = ?", *filters.Department)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, filters.SortOrder)
	}

	// Get records
	var dbEmployees []EmployeeModel
	err := query.Order(orderBy).
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&dbEmployees).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}

	employees := make([]*domain.Employee, len(dbEmployees))
	for i, dbEmployee := range dbEmployees {
		employees[i] = r.toDomainEmployee(&dbEmployee)
	}

	return employees, total, nil
}

// toDomainEmployee converts database model to domain entity
func (r *PostgresEmployeeRepository) toDomainEmployee(dbEmployee *EmployeeModel) *domain.Employee {
	employee := &domain.Employee{
		ID:           dbEmployee.ID,
		UserID:       dbEmployee.UserID,
		EmployeeCode: dbEmployee.EmployeeCode,
		Position:     dbEmployee.Position,
		Department:   dbEmployee.Department,
		HireDate:     dbEmployee.HireDate,
		Salary:       dbEmployee.Salary,
		HoursPerWeek: dbEmployee.HoursPerWeek,
		IsActive:     dbEmployee.IsActive,
		CreatedAt:    dbEmployee.CreatedAt,
		UpdatedAt:    dbEmployee.UpdatedAt,
	}

	// Convert user if loaded
	if dbEmployee.User != nil {
		employee.User = &domain.User{
			ID:        dbEmployee.User.ID,
			Email:     dbEmployee.User.Email,
			FirstName: dbEmployee.User.FirstName,
			LastName:  dbEmployee.User.LastName,
			Role:      dbEmployee.User.Role,
			IsActive:  dbEmployee.User.IsActive,
			CreatedAt: dbEmployee.User.CreatedAt,
			UpdatedAt: dbEmployee.User.UpdatedAt,
		}
	}

	return employee
}

// EmployeeModel represents the database table structure
type EmployeeModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	EmployeeCode string     `gorm:"uniqueIndex;not null"`
	Position     string     `gorm:"not null"`
	Department   string     `gorm:"not null;index"`
	HireDate     time.Time  `gorm:"not null"`
	Salary       float64    `gorm:"type:decimal(10,2)"`
	HoursPerWeek int        `gorm:"default:40"`
	IsActive     bool       `gorm:"default:true;index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `gorm:"index"`

	// Associations
	User *UserModel `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (EmployeeModel) TableName() string {
	return "employees"
}
