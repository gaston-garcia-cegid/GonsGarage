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

const sqlSelectCarBase = `SELECT id, make, model, year, license_plate, vin, color, mileage, client_id, created_at, updated_at, deleted_at
FROM cars WHERE deleted_at IS NULL`

// postgresCarRepository implements CarRepository using PostgreSQL
type postgresCarRepository struct {
	db   *gorm.DB
	sqlx *sqlx.DB
}

// NewPostgresCarRepository creates a new PostgreSQL car repository
func NewPostgresCarRepository(db *gorm.DB) ports.CarRepository {
	return &postgresCarRepository{db: db, sqlx: sqlxFromGORM(db)}
}

// CarModel represents the database table structure
type CarModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" db:"id"`
	Make         string     `gorm:"not null" db:"make"`
	Model        string     `gorm:"not null" db:"model"`
	Year         int        `gorm:"not null" db:"year"`
	LicensePlate string     `gorm:"column:license_plate;not null;uniqueIndex" db:"license_plate"`
	VIN          string     `gorm:"column:vin" db:"vin"`
	Color        string     `gorm:"not null" db:"color"`
	Mileage      int        `gorm:"not null;default:0" db:"mileage"`
	OwnerID      uuid.UUID  `gorm:"type:uuid;column:client_id;not null;index" db:"client_id"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" db:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index" db:"deleted_at"`

	// Relationships
	Owner   UserModel     `gorm:"foreignKey:OwnerID;references:ID"`
	Repairs []RepairModel `gorm:"foreignKey:CarID;references:ID"`
}

// TableName specifies the database table name
func (CarModel) TableName() string {
	return "cars"
}

// Create stores a new car in the database
func (r *postgresCarRepository) Create(ctx context.Context, car *domain.Car) error {
	if r.sqlx != nil {
		return r.createCarSQLX(ctx, car)
	}
	dbCar := r.toCarModel(car)
	if err := r.db.WithContext(ctx).Create(dbCar).Error; err != nil {
		return fmt.Errorf("failed to create car in database: %w", err)
	}
	return nil
}

func (r *postgresCarRepository) createCarSQLX(ctx context.Context, car *domain.Car) error {
	now := time.Now().UTC()
	const q = `INSERT INTO cars (id, make, model, year, license_plate, vin, color, mileage, client_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.sqlx.ExecContext(ctx, q,
		car.ID, car.Make, car.Model, car.Year, car.LicensePlate, car.VIN, car.Color, car.Mileage, car.OwnerID, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create car in database: %w", err)
	}
	car.CreatedAt = now
	car.UpdatedAt = now
	return nil
}

// GetByID retrieves a car by its unique identifier
func (r *postgresCarRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	if r.sqlx != nil {
		var row CarModel
		err := r.sqlx.GetContext(ctx, &row, sqlSelectCarBase+` AND id = $1`, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrCarNotFound
			}
			return nil, fmt.Errorf("failed to get car by ID: %w", err)
		}
		rows := []CarModel{row}
		if err := r.enrichCarsWithOwners(ctx, rows); err != nil {
			return nil, err
		}
		return r.toDomainCar(&rows[0]), nil
	}
	var dbCar CarModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Owner").
		First(&dbCar).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car by ID: %w", err)
	}
	return r.toDomainCar(&dbCar), nil
}

// GetByOwnerID retrieves all cars owned by a specific user
func (r *postgresCarRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Car, error) {
	if r.sqlx != nil {
		dbCars, err := r.selectCarsSQLX(ctx, "client_id = $1", []interface{}{ownerID}, 0, 0, "failed to get cars by owner ID")
		if err != nil {
			return nil, err
		}
		if err := r.enrichCarsWithOwners(ctx, dbCars); err != nil {
			return nil, err
		}
		return r.carsToDomain(dbCars), nil
	}
	var dbCars []CarModel
	err := r.db.WithContext(ctx).
		Where("client_id = ? AND deleted_at IS NULL", ownerID).
		Preload("Owner").
		Order("created_at DESC").
		Find(&dbCars).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get cars by owner ID: %w", err)
	}
	cars := make([]*domain.Car, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.toDomainCar(&dbCar)
	}
	return cars, nil
}

// GetByLicensePlate retrieves a car by its license plate
func (r *postgresCarRepository) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	if r.sqlx != nil {
		queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		var row CarModel
		err := r.sqlx.GetContext(queryCtx, &row, sqlSelectCarBase+` AND license_plate = $1`, licensePlate)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to get car by license plate: %w", err)
		}
		rows := []CarModel{row}
		if err := r.enrichCarsWithOwners(queryCtx, rows); err != nil {
			return nil, err
		}
		return r.toDomainCar(&rows[0]), nil
	}
	var dbCar CarModel
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := r.db.WithContext(queryCtx).
		Where("license_plate = ? AND deleted_at IS NULL", licensePlate).
		Preload("Owner").
		First(&dbCar).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get car by license plate: %w", err)
	}
	return r.toDomainCar(&dbCar), nil
}

// List retrieves cars with pagination
func (r *postgresCarRepository) List(ctx context.Context, limit, offset int) ([]*domain.Car, error) {
	if r.sqlx != nil {
		dbCars, err := r.selectCarsSQLX(ctx, "", nil, limit, offset, "failed to list cars")
		if err != nil {
			return nil, err
		}
		if err := r.enrichCarsWithOwners(ctx, dbCars); err != nil {
			return nil, err
		}
		return r.carsToDomain(dbCars), nil
	}
	var dbCars []CarModel
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Preload("Owner").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&dbCars).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list cars: %w", err)
	}
	cars := make([]*domain.Car, len(dbCars))
	for i, dbCar := range dbCars {
		cars[i] = r.toDomainCar(&dbCar)
	}
	return cars, nil
}

// Update modifies an existing car
func (r *postgresCarRepository) Update(ctx context.Context, car *domain.Car) error {
	if r.sqlx != nil {
		now := time.Now().UTC()
		const q = `UPDATE cars SET
make = $1, model = $2, year = $3, license_plate = $4, vin = $5, color = $6, mileage = $7, client_id = $8, updated_at = $9
WHERE id = $10 AND deleted_at IS NULL`
		res, err := r.sqlx.ExecContext(ctx, q,
			car.Make, car.Model, car.Year, car.LicensePlate, car.VIN, car.Color, car.Mileage, car.OwnerID, now, car.ID,
		)
		if err != nil {
			return fmt.Errorf("failed to update car: %w", err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read rows affected: %w", err)
		}
		if n == 0 {
			return domain.ErrCarNotFound
		}
		return nil
	}
	dbCar := r.toCarModel(car)
	result := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", car.ID).
		Updates(dbCar)
	if result.Error != nil {
		return fmt.Errorf("failed to update car: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrCarNotFound
	}
	return nil
}

// Delete removes a car using soft delete
func (r *postgresCarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if r.sqlx != nil {
		const q = `UPDATE cars SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
		res, err := r.sqlx.ExecContext(queryCtx, q, time.Now().UTC(), id)
		if err != nil {
			return fmt.Errorf("failed to delete car: %w", err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read rows affected: %w", err)
		}
		if n == 0 {
			return domain.ErrCarNotFound
		}
		return nil
	}
	result := r.db.WithContext(queryCtx).
		Model(&CarModel{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete car: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrCarNotFound
	}
	return nil
}

// GetWithRepairs retrieves a car with its complete repair history
func (r *postgresCarRepository) GetWithRepairs(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	if r.sqlx != nil {
		var row CarModel
		err := r.sqlx.GetContext(ctx, &row, sqlSelectCarBase+` AND id = $1`, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrCarNotFound
			}
			return nil, fmt.Errorf("failed to get car with repairs: %w", err)
		}
		rows := []CarModel{row}
		if err := r.enrichCarsWithOwners(ctx, rows); err != nil {
			return nil, err
		}
		dbCar := &rows[0]
		if err := r.db.WithContext(ctx).
			Where("car_id = ? AND deleted_at IS NULL", id).
			Order("created_at DESC").
			Find(&dbCar.Repairs).Error; err != nil {
			return nil, fmt.Errorf("failed to load repairs: %w", err)
		}
		return r.toDomainCarWithRepairs(dbCar), nil
	}
	var dbCar CarModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Owner").
		Preload("Repairs", "deleted_at IS NULL").
		First(&dbCar).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCarNotFound
		}
		return nil, fmt.Errorf("failed to get car with repairs: %w", err)
	}
	return r.toDomainCarWithRepairs(&dbCar), nil
}

// GetDeletedByLicensePlate retrieves a soft-deleted car by license plate
func (r *postgresCarRepository) GetDeletedByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if r.sqlx != nil {
		const q = `SELECT id, make, model, year, license_plate, vin, color, mileage, client_id, created_at, updated_at, deleted_at
FROM cars WHERE license_plate = $1 AND deleted_at IS NOT NULL LIMIT 1`
		var row CarModel
		err := r.sqlx.GetContext(queryCtx, &row, q, licensePlate)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to get deleted car: %w", err)
		}
		u, err := r.fetchUserByIDUnscoped(queryCtx, row.OwnerID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("failed to load car owner: %w", err)
			}
		} else {
			row.Owner = u
		}
		return r.toDomainCar(&row), nil
	}
	var dbCar CarModel
	err := r.db.WithContext(queryCtx).
		Unscoped().
		Where("license_plate = ? AND deleted_at IS NOT NULL", licensePlate).
		Preload("Owner").
		First(&dbCar).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get deleted car: %w", err)
	}
	return r.toDomainCar(&dbCar), nil
}

// Restore undeletes a soft-deleted car
func (r *postgresCarRepository) Restore(ctx context.Context, id uuid.UUID) error {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if r.sqlx != nil {
		const q = `UPDATE cars SET deleted_at = NULL, updated_at = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
		res, err := r.sqlx.ExecContext(queryCtx, q, time.Now().UTC(), id)
		if err != nil {
			return fmt.Errorf("failed to restore car: %w", err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read rows affected: %w", err)
		}
		if n == 0 {
			return domain.ErrCarNotFound
		}
		return nil
	}
	result := r.db.WithContext(queryCtx).
		Model(&CarModel{}).
		Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if result.Error != nil {
		return fmt.Errorf("failed to restore car: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrCarNotFound
	}
	return nil
}

func (r *postgresCarRepository) selectCarsSQLX(ctx context.Context, cond string, condArgs []interface{}, limit, offset int, errLabel string) ([]CarModel, error) {
	q := sqlSelectCarBase
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
	var rows []CarModel
	if err := r.sqlx.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("%s: %w", errLabel, err)
	}
	return rows, nil
}

func (r *postgresCarRepository) enrichCarsWithOwners(ctx context.Context, cars []CarModel) error {
	if len(cars) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, 0, len(cars))
	seen := make(map[uuid.UUID]struct{})
	for i := range cars {
		id := cars[i].OwnerID
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	m, err := r.fetchUsersByIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to load owners: %w", err)
	}
	for i := range cars {
		if u, ok := m[cars[i].OwnerID]; ok {
			cars[i].Owner = u
		}
	}
	return nil
}

func (r *postgresCarRepository) fetchUsersByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]UserModel, error) {
	return fetchUserModelsByIDs(ctx, r.sqlx, ids)
}

func (r *postgresCarRepository) fetchUserByIDUnscoped(ctx context.Context, id uuid.UUID) (UserModel, error) {
	const q = `SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at, deleted_at
FROM users WHERE id = $1`
	var u UserModel
	err := r.sqlx.GetContext(ctx, &u, q, id)
	return u, err
}

func (r *postgresCarRepository) carsToDomain(cars []CarModel) []*domain.Car {
	out := make([]*domain.Car, len(cars))
	for i := range cars {
		out[i] = r.toDomainCar(&cars[i])
	}
	return out
}

// Conversion methods following Clean Architecture principles

func (r *postgresCarRepository) toCarModel(car *domain.Car) *CarModel {
	return &CarModel{
		ID:           car.ID,
		Make:         car.Make,
		Model:        car.Model,
		Year:         car.Year,
		LicensePlate: car.LicensePlate,
		VIN:          car.VIN,
		Color:        car.Color,
		Mileage:      car.Mileage,
		OwnerID:      car.OwnerID,
		CreatedAt:    car.CreatedAt,
		UpdatedAt:    car.UpdatedAt,
		DeletedAt:    car.DeletedAt,
	}
}

func (r *postgresCarRepository) toDomainCar(dbCar *CarModel) *domain.Car {
	car := &domain.Car{
		ID:           dbCar.ID,
		Make:         dbCar.Make,
		Model:        dbCar.Model,
		Year:         dbCar.Year,
		LicensePlate: dbCar.LicensePlate,
		VIN:          dbCar.VIN,
		Color:        dbCar.Color,
		Mileage:      dbCar.Mileage,
		OwnerID:      dbCar.OwnerID,
		CreatedAt:    dbCar.CreatedAt,
		UpdatedAt:    dbCar.UpdatedAt,
		DeletedAt:    dbCar.DeletedAt,
	}

	// Convert owner if preloaded
	if dbCar.Owner.ID != uuid.Nil {
		car.Owner = r.toDomainUser(&dbCar.Owner)
	}

	return car
}

func (r *postgresCarRepository) toDomainCarWithRepairs(dbCar *CarModel) *domain.Car {
	car := r.toDomainCar(dbCar)

	// Convert repairs if preloaded
	if len(dbCar.Repairs) > 0 {
		car.Repairs = make([]domain.Repair, len(dbCar.Repairs))
		for i, dbRepair := range dbCar.Repairs {
			car.Repairs[i] = r.toDomainRepair(&dbRepair)
		}
	}

	return car
}

// Helper conversion methods (these would be shared or in a separate file)
func (r *postgresCarRepository) toDomainUser(dbUser *UserModel) domain.User {
	return domain.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Role:      dbUser.Role,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		DeletedAt: dbUser.DeletedAt,
	}
}

func (r *postgresCarRepository) toDomainRepair(dbRepair *RepairModel) domain.Repair {
	return domain.Repair{
		ID:           dbRepair.ID,
		CarID:        dbRepair.CarID,
		TechnicianID: dbRepair.TechnicianID,
		Description:  dbRepair.Description,
		Status:       domain.RepairStatus(dbRepair.Status),
		Cost:         dbRepair.Cost,
		StartedAt:    dbRepair.StartedAt,
		CompletedAt:  dbRepair.CompletedAt,
		CreatedAt:    dbRepair.CreatedAt,
		UpdatedAt:    dbRepair.UpdatedAt,
		DeletedAt:    dbRepair.DeletedAt,
	}
}

func (c *CarModel) toDomain() *domain.Car {
	car := &domain.Car{
		ID:           c.ID,
		Make:         c.Make,
		Model:        c.Model,
		Year:         c.Year,
		LicensePlate: c.LicensePlate,
		VIN:          c.VIN,
		Color:        c.Color,
		Mileage:      c.Mileage,
		OwnerID:      c.OwnerID,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		DeletedAt:    c.DeletedAt,
	}

	// Convert owner if present
	if c.Owner.ID != uuid.Nil {
		car.Owner = domain.User{
			ID:        c.Owner.ID,
			Email:     c.Owner.Email,
			FirstName: c.Owner.FirstName,
			LastName:  c.Owner.LastName,
			Role:      c.Owner.Role,
			IsActive:  c.Owner.IsActive,
			CreatedAt: c.Owner.CreatedAt,
			UpdatedAt: c.Owner.UpdatedAt,
			DeletedAt: c.Owner.DeletedAt,
		}
	}

	// Convert repairs if present (if RepairModel defined)
	if len(c.Repairs) > 0 {
		car.Repairs = make([]domain.Repair, len(c.Repairs))
		for i, r := range c.Repairs {
			car.Repairs[i] = domain.Repair{
				ID:           r.ID,
				CarID:        r.CarID,
				TechnicianID: r.TechnicianID,
				Description:  r.Description,
				Status:       domain.RepairStatus(r.Status),
				Cost:         r.Cost,
				StartedAt:    r.StartedAt,
				CompletedAt:  r.CompletedAt,
				CreatedAt:    r.CreatedAt,
				UpdatedAt:    r.UpdatedAt,
				DeletedAt:    r.DeletedAt,
			}
		}
	}

	return car
}
