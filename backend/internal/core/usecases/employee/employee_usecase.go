package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/backend/internal/core/ports"
	"github.com/google/uuid"
)

type EmployeeUseCase struct {
	employeeRepo ports.EmployeeRepository
	cache        ports.CacheRepository
}

func NewEmployeeUseCase(employeeRepo ports.EmployeeRepository, cache ports.CacheRepository) *EmployeeUseCase {
	return &EmployeeUseCase{
		employeeRepo: employeeRepo,
		cache:        cache,
	}
}

func (uc *EmployeeUseCase) CreateEmployee(ctx context.Context, req ports.CreateEmployeeRequest) (*domain.Employee, error) {
	employee := domain.NewEmployee(
		req.UserID,
		req.FirstName,
		req.LastName,
		req.Position,
		req.HourlyRate,
	)

	if err := uc.employeeRepo.Create(ctx, employee); err != nil {
		return nil, err
	}

	// Invalidate cache
	uc.cache.Delete(ctx, "employees:list")

	return employee, nil
}

func (uc *EmployeeUseCase) GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	cacheKey := fmt.Sprintf("employee:%s", id.String())

	// Try cache first
	var employee domain.Employee
	if err := uc.cache.Get(ctx, cacheKey, &employee); err == nil {
		return &employee, nil
	}

	// Get from database
	emp, err := uc.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Cache for 5 minutes
	uc.cache.Set(ctx, cacheKey, emp, 5*time.Minute)

	return emp, nil
}

func (uc *EmployeeUseCase) ListEmployees(ctx context.Context, limit, offset int) ([]*domain.Employee, error) {
	cacheKey := fmt.Sprintf("employees:list:%d:%d", limit, offset)

	// Try cache first
	var employees []*domain.Employee
	if err := uc.cache.Get(ctx, cacheKey, &employees); err == nil {
		return employees, nil
	}

	// Get from database
	employees, err := uc.employeeRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache for 1 minute
	uc.cache.Set(ctx, cacheKey, employees, 1*time.Minute)

	return employees, nil
}

func (uc *EmployeeUseCase) UpdateEmployee(ctx context.Context, id uuid.UUID, req ports.UpdateEmployeeRequest) (*domain.Employee, error) {
	employee, err := uc.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != nil {
		employee.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		employee.LastName = *req.LastName
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.HourlyRate != nil {
		employee.HourlyRate = *req.HourlyRate
	}
	if req.IsActive != nil {
		employee.IsActive = *req.IsActive
	}

	if err := uc.employeeRepo.Update(ctx, employee); err != nil {
		return nil, err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("employee:%s", id.String())
	uc.cache.Delete(ctx, cacheKey)
	uc.cache.Delete(ctx, "employees:list")

	return employee, nil
}

func (uc *EmployeeUseCase) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	if err := uc.employeeRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("employee:%s", id.String())
	uc.cache.Delete(ctx, cacheKey)
	uc.cache.Delete(ctx, "employees:list")

	return nil
}
