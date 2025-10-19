package domain

import (
	"time"

	"github.com/google/uuid"

	"errors"
)

type Employee struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	FirstName    string    `json:"first_name" gorm:"not null"`
	LastName     string    `json:"last_name" gorm:"not null"`
	Position     string    `json:"position" gorm:"not null"`
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
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func NewEmployee(firstName, lastName, position string, hourlyRate float64, role, department, phone string) *Employee {
	return &Employee{
		ID:         uuid.New(),
		FirstName:  firstName,
		LastName:   lastName,
		Position:   position,
		HourlyRate: hourlyRate,
		IsActive:   true,
		Role:       role,
		Department: department,
		Phone:      phone,
	}
}

func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

// ErrEmployeeNotFound is returned when an employee is not found in the repository.
var ErrEmployeeNotFound = errors.New("employee not found")
