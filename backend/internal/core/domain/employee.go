package domain

import (
	"time"

	"github.com/google/uuid"

	"errors"
)

type Employee struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserID       uuid.UUID  `json:"userID" gorm:"type:uuid;not null"`
	FirstName    string     `json:"firstName" gorm:"not null"`
	LastName     string     `json:"lastName" gorm:"not null"`
	Email        string     `json:"email" gorm:"unique;not null"`
	Phone        string     `json:"phone"`
	Department   string     `json:"department" gorm:"not null"`
	Position     string     `json:"position" gorm:"not null"`
	HourlyRate   float64    `json:"hourlyRate" gorm:"not null"`
	HoursWorked  float64    `json:"hoursWorked" gorm:"default:0"`
	IsActive     bool       `json:"isActive" gorm:"default:true"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	EmployeeCode string     `json:"employeeCode" gorm:"uniqueIndex;not null"`
	HireDate     time.Time  `json:"hireDate"`
	Salary       float64    `json:"salary"`
	HoursPerWeek int        `json:"hoursPerWeek"`
	Role         string     `json:"role"`
	PhoneNumber  string     `json:"phoneNumber"`
	DeletedAt    *time.Time `gorm:"index" json:"deletedAt,omitempty"`

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
