package domain

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	FirstName   string    `json:"first_name" gorm:"not null"`
	LastName    string    `json:"last_name" gorm:"not null"`
	Position    string    `json:"position" gorm:"not null"`
	HourlyRate  float64   `json:"hourly_rate" gorm:"not null"`
	HoursWorked float64   `json:"hours_worked" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func NewEmployee(userID uuid.UUID, firstName, lastName, position string, hourlyRate float64) *Employee {
	return &Employee{
		ID:         uuid.New(),
		UserID:     userID,
		FirstName:  firstName,
		LastName:   lastName,
		Position:   position,
		HourlyRate: hourlyRate,
		IsActive:   true,
	}
}

func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}
