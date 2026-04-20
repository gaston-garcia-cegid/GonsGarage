package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
)

// Columns match domain.Repair / GORM AutoMigrate (repairs has no denormalized car columns).
const sqlSelectRepairBase = `SELECT r.id, r.car_id, r.technician_id, r.description, r.status, r.cost, r.started_at, r.completed_at, r.created_at, r.updated_at, r.deleted_at
FROM repairs r WHERE r.deleted_at IS NULL`

type PostgresRepairRepository struct {
	db   *gorm.DB
	sqlx *sqlx.DB
}

func NewPostgresRepairRepository(db *gorm.DB) ports.RepairRepository {
	return &PostgresRepairRepository{db: db, sqlx: sqlxFromGORM(db)}
}

// RepairModel maps the repairs table as created by GORM from domain.Repair.
type RepairModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" db:"id"`
	CarID        uuid.UUID  `gorm:"type:uuid;not null;index" db:"car_id"`
	TechnicianID uuid.UUID  `gorm:"type:uuid;not null;index" db:"technician_id"`
	Description  string     `gorm:"not null" db:"description"`
	Status       string     `gorm:"not null;default:'pending'" db:"status"`
	Cost         float64    `gorm:"not null;default:0" db:"cost"`
	StartedAt    *time.Time `gorm:"column:started_at" db:"started_at"`
	CompletedAt  *time.Time `gorm:"column:completed_at" db:"completed_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index" db:"deleted_at"`
}

func (RepairModel) TableName() string {
	return "repairs"
}

func (r *PostgresRepairRepository) Create(ctx context.Context, repair *domain.Repair) error {
	if r.sqlx != nil {
		return r.createRepairSQLX(ctx, repair)
	}
	if err := r.db.WithContext(ctx).Create(repair).Error; err != nil {
		return fmt.Errorf("failed to create repair: %w", err)
	}
	return nil
}

func (r *PostgresRepairRepository) createRepairSQLX(ctx context.Context, repair *domain.Repair) error {
	started := time.Now().UTC()
	if repair.StartedAt != nil {
		started = repair.StartedAt.UTC()
	}
	const q = `INSERT INTO repairs (
id, car_id, technician_id, description, status, cost, started_at, completed_at, created_at, updated_at, deleted_at
) VALUES (
$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)`
	var completed interface{}
	if repair.CompletedAt != nil {
		completed = repair.CompletedAt.UTC()
	}
	_, err := r.sqlx.ExecContext(ctx, q,
		repair.ID,
		repair.CarID, repair.TechnicianID, repair.Description, string(repair.Status), repair.Cost,
		started, completed,
		repair.CreatedAt.UTC(), repair.UpdatedAt.UTC(), nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create repair: %w", err)
	}
	return nil
}

func (r *PostgresRepairRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error) {
	if r.sqlx != nil {
		var row RepairModel
		err := r.sqlx.GetContext(ctx, &row, sqlSelectRepairBase+` AND r.id = $1`, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrRepairNotFound
			}
			return nil, fmt.Errorf("failed to get repair by ID: %w", err)
		}
		return r.toDomainRepair(&row), nil
	}
	var dbRepair RepairModel
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&dbRepair).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRepairNotFound
		}
		return nil, fmt.Errorf("failed to get repair by ID: %w", err)
	}
	return r.toDomainRepair(&dbRepair), nil
}

func (r *PostgresRepairRepository) GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Repair, error) {
	if r.sqlx != nil {
		return r.selectRepairsSQLX(ctx, "r.car_id IN (SELECT id FROM cars WHERE owner_id = $1 AND deleted_at IS NULL)", []interface{}{clientID}, 0, 0, "failed to get repairs by client ID")
	}
	var dbRepairs []RepairModel
	err := r.db.WithContext(ctx).Table("repairs").
		Joins("JOIN cars ON cars.id = repairs.car_id").
		Where("cars.owner_id = ? AND repairs.deleted_at IS NULL AND cars.deleted_at IS NULL", clientID).
		Order("repairs.created_at DESC").Find(&dbRepairs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get repairs by client ID: %w", err)
	}
	return r.repairsToDomain(dbRepairs), nil
}

// GetByCarID implements ports.RepairRepository.
func (r *PostgresRepairRepository) GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error) {
	if r.sqlx != nil {
		return r.selectRepairsSQLX(ctx, "r.car_id = $1", []interface{}{carID}, 0, 0, "failed to get repairs by car ID")
	}
	var dbRepairs []RepairModel
	err := r.db.WithContext(ctx).Where("car_id = ? AND deleted_at IS NULL", carID).
		Order("created_at DESC").Find(&dbRepairs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get repairs by car ID: %w", err)
	}
	return r.repairsToDomain(dbRepairs), nil
}

func (r *PostgresRepairRepository) List(ctx context.Context, limit, offset int) ([]*domain.Repair, error) {
	if r.sqlx != nil {
		return r.selectRepairsSQLX(ctx, "", nil, limit, offset, "failed to list repairs")
	}
	var dbRepairs []RepairModel
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&dbRepairs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list repairs: %w", err)
	}
	return r.repairsToDomain(dbRepairs), nil
}

func (r *PostgresRepairRepository) selectRepairsSQLX(ctx context.Context, cond string, condArgs []interface{}, limit, offset int, errLabel string) ([]*domain.Repair, error) {
	q := sqlSelectRepairBase
	args := make([]interface{}, 0, 4+len(condArgs))
	if cond != "" {
		q += " AND " + cond
		args = append(args, condArgs...)
	}
	q += " ORDER BY created_at DESC"
	n := len(args)
	if limit > 0 {
		n++
		q += fmt.Sprintf(" LIMIT $%d", n)
		args = append(args, limit)
	}
	if offset > 0 {
		n++
		q += fmt.Sprintf(" OFFSET $%d", n)
		args = append(args, offset)
	}
	var rows []RepairModel
	if err := r.sqlx.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("%s: %w", errLabel, err)
	}
	return r.repairsToDomain(rows), nil
}

func (r *PostgresRepairRepository) repairsToDomain(rows []RepairModel) []*domain.Repair {
	out := make([]*domain.Repair, len(rows))
	for i := range rows {
		out[i] = r.toDomainRepair(&rows[i])
	}
	return out
}

func (r *PostgresRepairRepository) Update(ctx context.Context, repair *domain.Repair) error {
	if r.sqlx != nil {
		return r.updateRepairSQLX(ctx, repair)
	}
	dbRepair := &RepairModel{
		ID:        repair.ID,
		CreatedAt: repair.CreatedAt,
		UpdatedAt: repair.UpdatedAt,
	}
	result := r.db.WithContext(ctx).Model(dbRepair).Where("id = ?", repair.ID).Updates(dbRepair)
	if result.Error != nil {
		return fmt.Errorf("failed to update repair: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrRepairNotFound
	}
	return nil
}

func (r *PostgresRepairRepository) updateRepairSQLX(ctx context.Context, repair *domain.Repair) error {
	now := time.Now().UTC()
	started := now
	if repair.StartedAt != nil {
		started = repair.StartedAt.UTC()
	}
	const q = `UPDATE repairs SET
car_id = $1, technician_id = $2, description = $3, status = $4, cost = $5, started_at = $6, completed_at = $7, updated_at = $8
WHERE id = $9 AND deleted_at IS NULL`
	var completed interface{}
	if repair.CompletedAt != nil {
		completed = repair.CompletedAt.UTC()
	}
	res, err := r.sqlx.ExecContext(ctx, q,
		repair.CarID, repair.TechnicianID, repair.Description, string(repair.Status), repair.Cost,
		started, completed, now, repair.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update repair: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrRepairNotFound
	}
	return nil
}

func (r *PostgresRepairRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.sqlx != nil {
		const q = `UPDATE repairs SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
		res, err := r.sqlx.ExecContext(ctx, q, time.Now().UTC(), id)
		if err != nil {
			return fmt.Errorf("failed to delete repair: %w", err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read rows affected: %w", err)
		}
		if n == 0 {
			return domain.ErrRepairNotFound
		}
		return nil
	}
	result := r.db.WithContext(ctx).Model(&RepairModel{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete repair: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrRepairNotFound
	}
	return nil
}

func (r *PostgresRepairRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Repair, error) {
	if r.sqlx != nil {
		var row RepairModel
		err := r.sqlx.GetContext(ctx, &row, sqlSelectRepairBase+` AND EXISTS (SELECT 1 FROM cars c WHERE c.id = r.car_id AND c.license_plate = $1 AND c.deleted_at IS NULL)`, licensePlate)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrRepairNotFound
			}
			return nil, fmt.Errorf("failed to get repair by license plate: %w", err)
		}
		return r.toDomainRepair(&row), nil
	}
	var dbRepair RepairModel
	err := r.db.WithContext(ctx).Table("repairs").
		Joins("JOIN cars ON cars.id = repairs.car_id").
		Where("cars.license_plate = ? AND repairs.deleted_at IS NULL AND cars.deleted_at IS NULL", licensePlate).
		First(&dbRepair).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrRepairNotFound
		}
		return nil, fmt.Errorf("failed to get repair by license plate: %w", err)
	}
	return r.toDomainRepair(&dbRepair), nil
}

func (r *PostgresRepairRepository) toDomainRepair(dbRepair *RepairModel) *domain.Repair {
	var started *time.Time
	if dbRepair.StartedAt != nil && !dbRepair.StartedAt.IsZero() {
		t := *dbRepair.StartedAt
		started = &t
	}
	return &domain.Repair{
		ID:           dbRepair.ID,
		CarID:        dbRepair.CarID,
		TechnicianID: dbRepair.TechnicianID,
		Description:  dbRepair.Description,
		Status:       domain.RepairStatus(dbRepair.Status),
		Cost:         dbRepair.Cost,
		StartedAt:    started,
		CompletedAt:  dbRepair.CompletedAt,
		CreatedAt:    dbRepair.CreatedAt,
		UpdatedAt:    dbRepair.UpdatedAt,
		DeletedAt:    dbRepair.DeletedAt,
	}
}
