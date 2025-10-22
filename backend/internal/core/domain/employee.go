package domain

import (
	"time"

	"github.com/google/uuid"

	"errors"
)

type Employee struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	FirstName    string     `json:"first_name" gorm:"not null"`
	LastName     string     `json:"last_name" gorm:"not null"`
	Email        string     `json:"email" gorm:"unique;not null"`
	Phone        string     `json:"phone"`
	Department   string     `json:"department" gorm:"not null"`
	Position     string     `json:"position" gorm:"not null"`
	HourlyRate   float64    `json:"hourly_rate" gorm:"not null"`
	HoursWorked  float64    `json:"hours_worked" gorm:"default:0"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	EmployeeCode string     `json:"employee_code"`
	HireDate     time.Time  `json:"hire_date"`
	Salary       float64    `json:"salary"`
	HoursPerWeek int        `json:"hours_per_week"`
	Role         string     `json:"role"`
	PhoneNumber  string     `json:"phone_number"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for GORM
func (Employee) TableName() string {
	return "employees"
}

func NewEmployee(firstName, lastName, email, position string, hourlyRate float64, role, department, phone string, hireDate time.Time, salary float64, hoursPerWeek int) *Employee {
	return &Employee{
		ID:           uuid.New(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Position:     position,
		HireDate:     time.Now(),
		Salary:       salary,
		HourlyRate:   hourlyRate,
		IsActive:     true,
		Role:         role,
		Department:   department,
		Phone:        phone,
		EmployeeCode: generateEmployeeCode(),
		HoursPerWeek: 40,
	}
}

func generateEmployeeCode() string {
	// Gerar código único (pode ser melhorado)
	return uuid.New().String()[:8]
}

func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

// ErrEmployeeNotFound is returned when an employee is not found in the repository.
var ErrEmployeeNotFound = errors.New("employee not found")
