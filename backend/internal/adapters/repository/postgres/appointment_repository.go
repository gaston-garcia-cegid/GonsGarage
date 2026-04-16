package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// postgresAppointmentRepository implements AppointmentRepository using PostgreSQL
type postgresAppointmentRepository struct {
	db *gorm.DB
}

// NewPostgresAppointmentRepository creates a new PostgreSQL appointment repository
func NewPostgresAppointmentRepository(db *gorm.DB) ports.AppointmentRepository {
	return &postgresAppointmentRepository{db: db}
}

// AppointmentModel represents the database table structure
type AppointmentModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID    uuid.UUID  `gorm:"type:uuid"`
	CarID         uuid.UUID  `gorm:"type:uuid"`
	ScheduledTime time.Time  `gorm:"column:scheduled_at;type:timestamptz"`
	Notes         string     `gorm:"column:notes;type:text"`
	Status        string     `gorm:"type:text"`
	ServiceType   string     `gorm:"column:service_type;type:text"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index"`
}

// TableName specifies the database table name
func (AppointmentModel) TableName() string {
	return "appointments"
}

// Create stores a new appointment in the database
func (r *postgresAppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	dbAppointment := r.toAppointmentModel(appointment)
	if err := r.db.WithContext(ctx).Create(&dbAppointment).Error; err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	return nil
}

// toAppointmentModel converts a domain.Appointment to AppointmentModel
func (r *postgresAppointmentRepository) toAppointmentModel(appointment *domain.Appointment) *AppointmentModel {
	return &AppointmentModel{
		ID:            appointment.ID,
		CustomerID:    appointment.CustomerID,
		CarID:         appointment.CarID,
		ScheduledTime: appointment.ScheduledAt,
		Notes:         appointment.Notes,
		Status:        string(appointment.Status),
		ServiceType:   appointment.ServiceType,
		CreatedAt:     appointment.CreatedAt,
		UpdatedAt:     appointment.UpdatedAt,
		DeletedAt:     appointment.DeletedAt,
	}
}

// toDomainAppointment converts an AppointmentModel to domain.Appointment
func (r *postgresAppointmentRepository) toDomainAppointment(model *AppointmentModel) *domain.Appointment {
	return &domain.Appointment{
		ID:          model.ID,
		CustomerID:  model.CustomerID,
		CarID:       model.CarID,
		ScheduledAt: model.ScheduledTime,
		Notes:       model.Notes,
		Status:      domain.AppointmentStatus(model.Status),
		ServiceType: model.ServiceType,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   model.DeletedAt,
	}
}

// GetByID retrieves an appointment by its ID
func (r *postgresAppointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	var model AppointmentModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrAppointmentNotFound
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}
	return r.toDomainAppointment(&model), nil
}

// Update modifies an existing appointment in the database
func (r *postgresAppointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	dbAppointment := r.toAppointmentModel(appointment)
	if err := r.db.WithContext(ctx).Save(dbAppointment).Error; err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}
	return nil
}

// Delete removes an appointment from the database
func (r *postgresAppointmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&AppointmentModel{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}
	return nil
}

// List retrieves appointments with optional filters
func (r *postgresAppointmentRepository) List(ctx context.Context, filters *ports.AppointmentFilters) ([]*domain.Appointment, int64, error) {
	buildQuery := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&AppointmentModel{}).Where("deleted_at IS NULL")
		if filters == nil {
			return q
		}
		if filters.CustomerID != nil {
			q = q.Where("customer_id = ?", *filters.CustomerID)
		}
		if filters.CarID != nil {
			q = q.Where("car_id = ?", *filters.CarID)
		}
		if filters.Status != nil && *filters.Status != "" {
			q = q.Where("status = ?", *filters.Status)
		}
		return q
	}

	var total int64
	if err := buildQuery().Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count appointments: %w", err)
	}

	orderCol := "created_at"
	if filters != nil && filters.SortBy != "" {
		switch filters.SortBy {
		case "scheduled_at":
			orderCol = "scheduled_at"
		case "created_at":
			orderCol = "created_at"
		default:
			orderCol = "created_at"
		}
	}
	dir := "DESC"
	if filters != nil && strings.ToUpper(filters.SortOrder) == "ASC" {
		dir = "ASC"
	}

	limit := 10
	offset := 0
	if filters != nil {
		if filters.Limit > 0 {
			limit = filters.Limit
		}
		offset = filters.Offset
	}

	var models []AppointmentModel
	if err := buildQuery().
		Order(orderCol + " " + dir).
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list appointments: %w", err)
	}

	appointments := make([]*domain.Appointment, 0, len(models))
	for i := range models {
		appointments = append(appointments, r.toDomainAppointment(&models[i]))
	}
	return appointments, total, nil
}
