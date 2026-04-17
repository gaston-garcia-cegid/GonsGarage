package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

const sqlAppointmentSelectList = `SELECT id, customer_id, car_id, scheduled_at, notes, status, service_type, created_at, updated_at, deleted_at FROM appointments`

// postgresAppointmentRepository implements AppointmentRepository using PostgreSQL
type postgresAppointmentRepository struct {
	db   *gorm.DB
	sqlx *sqlx.DB
}

// NewPostgresAppointmentRepository creates a new PostgreSQL appointment repository
func NewPostgresAppointmentRepository(db *gorm.DB) ports.AppointmentRepository {
	return &postgresAppointmentRepository{db: db, sqlx: sqlxFromGORM(db)}
}

// AppointmentModel represents the database table structure
type AppointmentModel struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" db:"id"`
	CustomerID    uuid.UUID  `gorm:"type:uuid" db:"customer_id"`
	CarID         uuid.UUID  `gorm:"type:uuid" db:"car_id"`
	ScheduledTime time.Time  `gorm:"column:scheduled_at;type:timestamptz" db:"scheduled_at"`
	Notes         string     `gorm:"column:notes;type:text" db:"notes"`
	Status        string     `gorm:"type:text" db:"status"`
	ServiceType   string     `gorm:"column:service_type;type:text" db:"service_type"`
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at;index" db:"deleted_at"`
}

// TableName specifies the database table name
func (AppointmentModel) TableName() string {
	return "appointments"
}

// Create stores a new appointment in the database
func (r *postgresAppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	if r.sqlx != nil {
		return r.createAppointmentSQLX(ctx, appointment)
	}
	dbAppointment := r.toAppointmentModel(appointment)
	if err := r.db.WithContext(ctx).Create(&dbAppointment).Error; err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	return nil
}

func (r *postgresAppointmentRepository) createAppointmentSQLX(ctx context.Context, appointment *domain.Appointment) error {
	now := time.Now().UTC()
	const q = `INSERT INTO appointments (id, customer_id, car_id, scheduled_at, notes, status, service_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.sqlx.ExecContext(ctx, q,
		appointment.ID, appointment.CustomerID, appointment.CarID, appointment.ScheduledAt.UTC(),
		appointment.Notes, string(appointment.Status), appointment.ServiceType,
		now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create appointment: %w", err)
	}
	appointment.CreatedAt = now
	appointment.UpdatedAt = now
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
	if r.sqlx != nil {
		var model AppointmentModel
		err := r.sqlx.GetContext(ctx, &model, sqlAppointmentSelectList+` WHERE deleted_at IS NULL AND id = $1`, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrAppointmentNotFound
			}
			return nil, fmt.Errorf("failed to get appointment: %w", err)
		}
		return r.toDomainAppointment(&model), nil
	}
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
	if r.sqlx != nil {
		return r.updateAppointmentSQLX(ctx, appointment)
	}
	dbAppointment := r.toAppointmentModel(appointment)
	if err := r.db.WithContext(ctx).Save(dbAppointment).Error; err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}
	return nil
}

func (r *postgresAppointmentRepository) updateAppointmentSQLX(ctx context.Context, appointment *domain.Appointment) error {
	now := time.Now().UTC()
	const q = `UPDATE appointments SET
customer_id = $1, car_id = $2, scheduled_at = $3, notes = $4, status = $5, service_type = $6, updated_at = $7
WHERE id = $8 AND deleted_at IS NULL`
	res, err := r.sqlx.ExecContext(ctx, q,
		appointment.CustomerID, appointment.CarID, appointment.ScheduledAt.UTC(),
		appointment.Notes, string(appointment.Status), appointment.ServiceType, now, appointment.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrAppointmentNotFound
	}
	appointment.UpdatedAt = now
	return nil
}

// Delete removes an appointment from the database (soft delete)
func (r *postgresAppointmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.sqlx != nil {
		const q = `UPDATE appointments SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
		_, err := r.sqlx.ExecContext(ctx, q, time.Now().UTC(), id)
		if err != nil {
			return fmt.Errorf("failed to delete appointment: %w", err)
		}
		return nil
	}
	if err := r.db.WithContext(ctx).Delete(&AppointmentModel{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}
	return nil
}

// List retrieves appointments with optional filters
func (r *postgresAppointmentRepository) List(ctx context.Context, filters *ports.AppointmentFilters) ([]*domain.Appointment, int64, error) {
	if r.sqlx != nil {
		return r.listAppointmentsSQLX(ctx, filters)
	}
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

func (r *postgresAppointmentRepository) listAppointmentsSQLX(ctx context.Context, filters *ports.AppointmentFilters) ([]*domain.Appointment, int64, error) {
	where := []string{"deleted_at IS NULL"}
	args := make([]interface{}, 0, 8)
	if filters != nil {
		if filters.CustomerID != nil {
			where = append(where, fmt.Sprintf("customer_id = $%d", len(args)+1))
			args = append(args, *filters.CustomerID)
		}
		if filters.CarID != nil {
			where = append(where, fmt.Sprintf("car_id = $%d", len(args)+1))
			args = append(args, *filters.CarID)
		}
		if filters.Status != nil && *filters.Status != "" {
			where = append(where, fmt.Sprintf("status = $%d", len(args)+1))
			args = append(args, *filters.Status)
		}
	}
	whereSQL := strings.Join(where, " AND ")

	countQ := `SELECT COUNT(*) FROM appointments WHERE ` + whereSQL
	var total int64
	if err := r.sqlx.GetContext(ctx, &total, countQ, args...); err != nil {
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

	listArgs := append(append([]interface{}(nil), args...), limit, offset)
	n := len(args)
	listQ := sqlAppointmentSelectList + ` WHERE ` + whereSQL + fmt.Sprintf(` ORDER BY %s %s LIMIT $%d OFFSET $%d`, orderCol, dir, n+1, n+2)

	var models []AppointmentModel
	if err := r.sqlx.SelectContext(ctx, &models, listQ, listArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to list appointments: %w", err)
	}

	out := make([]*domain.Appointment, 0, len(models))
	for i := range models {
		out = append(out, r.toDomainAppointment(&models[i]))
	}
	return out, total, nil
}
