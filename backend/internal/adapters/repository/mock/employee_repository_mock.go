package mock

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

type EmployeeRepositoryMock struct {
	employees map[uuid.UUID]*domain.Employee
}

func NewEmployeeRepositoryMock() *EmployeeRepositoryMock {
	return &EmployeeRepositoryMock{
		employees: make(map[uuid.UUID]*domain.Employee),
	}
}

// Implement the methods of the EmployeeRepository interface
func (m *EmployeeRepositoryMock) Create(ctx context.Context, employee *domain.Employee) error {
	m.employees[employee.ID] = employee
	return nil
}
func (m *EmployeeRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	employee, exists := m.employees[id]
	if !exists {
		return nil, domain.ErrEmployeeNotFound
	}
	return employee, nil
}
func (m *EmployeeRepositoryMock) Update(ctx context.Context, employee *domain.Employee) error {
	m.employees[employee.ID] = employee
	return nil
}
func (m *EmployeeRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	delete(m.employees, id)
	return nil
}
func (m *EmployeeRepositoryMock) GetByDepartment(ctx context.Context, department string) ([]*domain.Employee, error) {
	var result []*domain.Employee
	for _, employee := range m.employees {
		if employee.Department == department {
			result = append(result, employee)
		}
	}
	return result, nil
}
func (m *EmployeeRepositoryMock) GetActiveEmployees(ctx context.Context) ([]*domain.Employee, error) {
	var result []*domain.Employee
	for _, employee := range m.employees {
		if employee.IsActive {
			result = append(result, employee)
		}
	}
	return result, nil
}
func (m *EmployeeRepositoryMock) GetEmployeesByRole(ctx context.Context, role string) ([]*domain.Employee, error) {
	var result []*domain.Employee
	for _, employee := range m.employees {
		if employee.Role == role {
			result = append(result, employee)
		}
	}
	return result, nil
}
func (m *EmployeeRepositoryMock) List(ctx context.Context, filters *domain.EmployeeFilters) ([]*domain.Employee, int64, error) {
	var result []*domain.Employee
	for _, employee := range m.employees {
		if filters.Department != "" && employee.Department != filters.Department {
			continue
		}
		if filters.Role != "" && employee.Role != filters.Role {
			continue
		}
		result = append(result, employee)
	}
	return result, int64(len(result)), nil
}
func (m *EmployeeRepositoryMock) GetEmployeeWorkHours(ctx context.Context, employeeID string, from, to time.Time) ([]*domain.WorkHour, error) {
	// Mock implementation, return empty slice
	return []*domain.WorkHour{}, nil
}
func (m *EmployeeRepositoryMock) GetEmployeesOnVacation(ctx context.Context, date time.Time) ([]*domain.Employee, error) {
	// Mock implementation, return empty slice
	return []*domain.Employee{}, nil
}
func (m *EmployeeRepositoryMock) SearchByName(ctx context.Context, name string, limit int) ([]*domain.Employee, error) {
	var result []*domain.Employee
	for _, employee := range m.employees {
		if employee.FirstName == name {
			result = append(result, employee)
		}
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}
