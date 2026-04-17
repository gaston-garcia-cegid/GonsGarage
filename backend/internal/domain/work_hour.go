package domain

import (
	"time"

	"github.com/google/uuid"
)

// WorkHour representa as horas trabalhadas de um funcionário
type WorkHour struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	EmployeeID  uuid.UUID `json:"employee_id" gorm:"type:uuid;not null"`
	Date        time.Time `json:"date" gorm:"type:date;not null"`
	StartTime   time.Time `json:"start_time" gorm:"type:time;not null"`
	EndTime     time.Time `json:"end_time" gorm:"type:time"`
	Hours       float64   `json:"hours" gorm:"type:decimal(5,2)"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relacionamento
	Employee *Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

// NewWorkHour cria um novo registro de horas trabalhadas
func NewWorkHour(employeeID uuid.UUID, date, startTime time.Time) *WorkHour {
	return &WorkHour{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		Date:       date,
		StartTime:  startTime,
	}
}

// CalculateHours calcula as horas trabalhadas
func (w *WorkHour) CalculateHours() float64 {
	if w.EndTime.IsZero() {
		return 0
	}
	duration := w.EndTime.Sub(w.StartTime)
	return duration.Hours()
}
