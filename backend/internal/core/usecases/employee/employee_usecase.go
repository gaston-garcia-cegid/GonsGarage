package employee

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
)

type EmployeeUseCase struct {
	employeeRepo ports.EmployeeRepository
	cache        ports.CacheRepository
}

func NewEmployeeUseCase(
	employeeRepo ports.EmployeeRepository,
	cache ports.CacheRepository,
) ports.EmployeeService {
	return &EmployeeUseCase{
		employeeRepo: employeeRepo,
		cache:        cache,
	}
}

func (uc *EmployeeUseCase) CreateEmployee(ctx context.Context, req ports.CreateEmployeeRequest) (*domain.Employee, error) {
	employee := domain.NewEmployee(
		req.FirstName,
		req.LastName,
		req.Position,
		req.HourlyRate,
		req.Role,
		req.Department,
		req.PhoneNumber,
	)

	if err := uc.employeeRepo.Create(ctx, employee); err != nil {
		return nil, err
	}

	// Invalidate cache
	uc.cache.Delete(ctx, "employees:list")

	return employee, nil
}

func (uc *EmployeeUseCase) GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	return uc.employeeRepo.FindByID(ctx, uint(id.ID()))
}

func (uc *EmployeeUseCase) ListEmployees(ctx context.Context, limit, offset int) ([]*domain.Employee, error) {
	employees, total, err := uc.employeeRepo.List(ctx, limit, offset)
	_ = total // If you need total, you can return it or use it as needed
	return employees, err
}

func (uc *EmployeeUseCase) UpdateEmployee(ctx context.Context, id uuid.UUID, req ports.UpdateEmployeeRequest) (*domain.Employee, error) {
	employee, err := uc.employeeRepo.FindByID(ctx, uint(id.ID()))
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		employee.FirstName = *req.Name
	}
	if req.Email != nil {
		employee.User.Email = *req.Email
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.Department != nil {
		employee.Department = *req.Department
	}

	if err := uc.employeeRepo.Update(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (uc *EmployeeUseCase) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return uc.employeeRepo.Delete(ctx, uint(id.ID()))
}
