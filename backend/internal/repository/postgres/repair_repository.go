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

// Columns match GORM/AutoMigrate naming for repairs (includes denormalized car fields on the model).
const sqlSelectRepairBase = `SELECT id, client_id, make, model, year, license_plate, vin, color, mileage, created_at, updated_at, deleted_at, car_id, technician_id, description, status, cost, start_date, completed_at
FROM repairs WHERE deleted_at IS NULL`

type PostgresRepairRepository struct {
	db   *gorm.DB
	sqlx *sqlx.DB
}

func NewPostgresRepairRepository(db *gorm.DB) ports.RepairRepository {
	return &PostgresRepairRepository{db: db, sqlx: sqlxFromGORM(db)}
}

// RepairModel represents the database table structure
type RepairModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" db:"id"`
	ClientID     uuid.UUID  `gorm:"type:uuid;not null;index" db:"client_id"`
	Make         string     `gorm:"not null" db:"make"`
	Model        string     `gorm:"not null" db:"model"`
	Year         int        `gorm:"not null" db:"year"`
	LicensePlate string     `gorm:"column:license_plate;not null;uniqueIndex" db:"license_plate"`
	VIN          string     `gorm:"column:vin" db:"vin"`
	Color        string     `gorm:"not null" db:"color"`
	Mileage      int        `gorm:"not null;default:0" db:"mileage"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index" db:"deleted_at"`
	CarID        uuid.UUID  `gorm:"type:uuid;not null;index" db:"car_id"`
	TechnicianID uuid.UUID  `gorm:"type:uuid;not null;index" db:"technician_id"`
	Description  string     `gorm:"not null" db:"description"`
	Status       string     `gorm:"not null;default:'pending'" db:"status"`
	Cost         float64    `gorm:"not null;default:0" db:"cost"`
	StartedAt    time.Time  `gorm:"column:start_date" db:"start_date"`
	CompletedAt  *time.Time `gorm:"column:completed_at" db:"completed_at"`
}

func (RepairModel) TableName() string {
	return "repairs"
}

func (r *PostgresRepairRepository) Create(ctx context.Context, repair *domain.Repair) error {
	if r.sqlx != nil {
		return r.createRepairSQLX(ctx, repair)
	}
	dbRepair := &RepairModel{
		ID:        repair.ID,
		CreatedAt: repair.CreatedAt,
		UpdatedAt: repair.UpdatedAt,
	}
	if err := r.db.WithContext(ctx).Create(dbRepair).Error; err != nil {
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
id, client_id, make, model, year, license_plate, vin, color, mileage,
created_at, updated_at, deleted_at,
car_id, technician_id, description, status, cost, start_date, completed_at
) VALUES (
$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
)`
	var completed interface{}
	if repair.CompletedAt != nil {
		completed = repair.CompletedAt.UTC()
	}
	_, err := r.sqlx.ExecContext(ctx, q,
		repair.ID,
		uuid.Nil, "", "", 0, "", "", "", 0,
		repair.CreatedAt.UTC(), repair.UpdatedAt.UTC(), nil,
		repair.CarID, repair.TechnicianID, repair.Description, string(repair.Status), repair.Cost,
		started, completed,
	)
	if err != nil {
		return fmt.Errorf("failed to create repair: %w", err)
	}
	return nil
}

func (r *PostgresRepairRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error) {
	if r.sqlx != nil {
		var row RepairModel
		err := r.sqlx.GetContext(ctx, &row, sqlSelectRepairBase+` AND id = $1`, id)
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
		return r.selectRepairsSQLX(ctx, "client_id = $1", []interface{}{clientID}, 0, 0, "failed to get repairs by client ID")
	}
	var dbRepairs []RepairModel
	err := r.db.WithContext(ctx).Where("client_id = ? AND deleted_at IS NULL", clientID).
		Order("created_at DESC").Find(&dbRepairs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get repairs by client ID: %w", err)
	}
	return r.repairsToDomain(dbRepairs), nil
}

// GetByCarID implements ports.RepairRepository.
func (r *PostgresRepairRepository) GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error) {
	if r.sqlx != nil {
		return r.selectRepairsSQLX(ctx, "car_id = $1", []interface{}{carID}, 0, 0, "failed to get repairs by car ID")
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
car_id = $1, technician_id = $2, description = $3, status = $4, cost = $5, start_date = $6, completed_at = $7, updated_at = $8
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
		err := r.sqlx.GetContext(ctx, &row, sqlSelectRepairBase+` AND license_plate = $1`, licensePlate)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrRepairNotFound
			}
			return nil, fmt.Errorf("failed to get repair by license plate: %w", err)
		}
		return r.toDomainRepair(&row), nil
	}
	var dbRepair RepairModel
	err := r.db.WithContext(ctx).Where("license_plate = ? AND deleted_at IS NULL", licensePlate).First(&dbRepair).Error
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
	if !dbRepair.StartedAt.IsZero() {
		t := dbRepair.StartedAt
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
