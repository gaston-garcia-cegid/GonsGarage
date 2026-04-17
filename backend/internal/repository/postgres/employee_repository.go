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

const sqlSelectEmployeeBase = `SELECT id, user_id, employee_code, position, department, hire_date, salary, hours_per_week, is_active, created_at, updated_at, deleted_at
FROM employees WHERE deleted_at IS NULL`

// PostgresEmployeeRepository implements EmployeeRepository interface
type PostgresEmployeeRepository struct {
	db   *gorm.DB
	sqlx *sqlx.DB
}

// NewPostgresEmployeeRepository creates a new PostgreSQL employee repository
func NewPostgresEmployeeRepository(db *gorm.DB) ports.EmployeeRepository {
	return &PostgresEmployeeRepository{db: db, sqlx: sqlxFromGORM(db)}
}

// Create implements EmployeeRepository.Create
func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	if r.sqlx != nil {
		return r.createEmployeeSQLX(ctx, employee)
	}
	dbEmployee := &EmployeeModel{
		ID:           employee.ID,
		UserID:       employee.UserID,
		EmployeeCode: employee.EmployeeCode,
		Position:     employee.Position,
		Department:   employee.Department,
		HireDate:     employee.HireDate,
		Salary:       employee.Salary,
		HoursPerWeek: employee.HoursPerWeek,
		IsActive:     employee.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(dbEmployee).Error; err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	employee.CreatedAt = dbEmployee.CreatedAt
	employee.UpdatedAt = dbEmployee.UpdatedAt
	return nil
}

func (r *PostgresEmployeeRepository) createEmployeeSQLX(ctx context.Context, employee *domain.Employee) error {
	now := time.Now().UTC()
	const q = `INSERT INTO employees (id, user_id, employee_code, position, department, hire_date, salary, hours_per_week, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.sqlx.ExecContext(ctx, q,
		employee.ID, employee.UserID, employee.EmployeeCode, employee.Position, employee.Department,
		employee.HireDate, employee.Salary, employee.HoursPerWeek, employee.IsActive, now, now,
	)
	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	employee.CreatedAt = now
	employee.UpdatedAt = now
	return nil
}

// GetByID loads an employee by UUID string (non-interface helper).
func (r *PostgresEmployeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	uid, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, fmt.Errorf("invalid employee id: %w", err)
	}
	return r.FindByID(ctx, uid)
}

// GetByDepartment returns employees in a department (extension beyond ports.EmployeeRepository).
func (r *PostgresEmployeeRepository) GetByDepartment(ctx context.Context, department string) ([]*domain.Employee, error) {
	if r.sqlx != nil {
		rows, err := r.selectEmployeesSQLX(ctx, "department = $1", []interface{}{department}, 0, 0, "failed to get employees by department")
		if err != nil {
			return nil, err
		}
		if err := r.enrichEmployeesWithUsers(ctx, rows); err != nil {
			return nil, err
		}
		return r.employeesToDomain(rows), nil
	}
	var dbEmployees []EmployeeModel
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("department = ? AND deleted_at IS NULL", department).
		Find(&dbEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by department: %w", err)
	}
	return r.employeesToDomain(dbEmployees), nil
}

// GetActiveEmployees returns active employees (extension beyond ports.EmployeeRepository).
func (r *PostgresEmployeeRepository) GetActiveEmployees(ctx context.Context) ([]*domain.Employee, error) {
	if r.sqlx != nil {
		rows, err := r.selectEmployeesSQLX(ctx, "is_active = $1", []interface{}{true}, 0, 0, "failed to get active employees")
		if err != nil {
			return nil, err
		}
		if err := r.enrichEmployeesWithUsers(ctx, rows); err != nil {
			return nil, err
		}
		return r.employeesToDomain(rows), nil
	}
	var dbEmployees []EmployeeModel
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("is_active = ? AND deleted_at IS NULL", true).
		Find(&dbEmployees).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active employees: %w", err)
	}
	return r.employeesToDomain(dbEmployees), nil
}

// Update implements EmployeeRepository.Update
func (r *PostgresEmployeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	if r.sqlx != nil {
		return r.updateEmployeeSQLX(ctx, employee)
	}
	dbEmployee := &EmployeeModel{
		ID:           employee.ID,
		UserID:       employee.UserID,
		EmployeeCode: employee.EmployeeCode,
		Position:     employee.Position,
		Department:   employee.Department,
		HireDate:     employee.HireDate,
		Salary:       employee.Salary,
		HoursPerWeek: employee.HoursPerWeek,
		IsActive:     employee.IsActive,
		UpdatedAt:    time.Now(),
	}
	result := r.db.WithContext(ctx).Model(dbEmployee).Where("id = ?", employee.ID).Updates(dbEmployee)
	if result.Error != nil {
		return fmt.Errorf("failed to update employee: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEmployeeNotFound
	}
	employee.UpdatedAt = dbEmployee.UpdatedAt
	return nil
}

func (r *PostgresEmployeeRepository) updateEmployeeSQLX(ctx context.Context, employee *domain.Employee) error {
	now := time.Now().UTC()
	const q = `UPDATE employees SET
user_id = $1, employee_code = $2, position = $3, department = $4, hire_date = $5, salary = $6, hours_per_week = $7, is_active = $8, updated_at = $9
WHERE id = $10 AND deleted_at IS NULL`
	res, err := r.sqlx.ExecContext(ctx, q,
		employee.UserID, employee.EmployeeCode, employee.Position, employee.Department, employee.HireDate,
		employee.Salary, employee.HoursPerWeek, employee.IsActive, now, employee.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrEmployeeNotFound
	}
	employee.UpdatedAt = now
	return nil
}

// Delete implements EmployeeRepository.Delete (soft delete)
func (r *PostgresEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.sqlx != nil {
		const q = `UPDATE employees SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
		res, err := r.sqlx.ExecContext(ctx, q, time.Now().UTC(), id)
		if err != nil {
			return fmt.Errorf("failed to delete employee: %w", err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read rows affected: %w", err)
		}
		if n == 0 {
			return domain.ErrEmployeeNotFound
		}
		return nil
	}
	result := r.db.WithContext(ctx).Model(&EmployeeModel{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrEmployeeNotFound
	}
	return nil
}

// List implements EmployeeRepository.List
func (r *PostgresEmployeeRepository) List(ctx context.Context, filters *ports.EmployeeFilters) ([]*domain.Employee, int64, error) {
	if filters == nil {
		filters = &ports.EmployeeFilters{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"}
	}
	if r.sqlx != nil {
		return r.listEmployeesSQLX(ctx, filters)
	}
	query := r.db.WithContext(ctx).Model(&EmployeeModel{}).
		Preload("User").
		Where("employees.deleted_at IS NULL")
	if filters.Department != nil {
		query = query.Where("department = ?", *filters.Department)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}
	orderBy := "created_at DESC"
	if filters.SortBy != "" {
		orderBy = fmt.Sprintf("%s %s", filters.SortBy, filters.SortOrder)
	}
	var dbEmployees []EmployeeModel
	err := query.Order(orderBy).
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&dbEmployees).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}
	return r.employeesToDomain(dbEmployees), total, nil
}

func (r *PostgresEmployeeRepository) listEmployeesSQLX(ctx context.Context, filters *ports.EmployeeFilters) ([]*domain.Employee, int64, error) {
	where := []string{"deleted_at IS NULL"}
	args := make([]interface{}, 0, 6)
	if filters != nil {
		if filters.Department != nil {
			where = append(where, fmt.Sprintf("department = $%d", len(args)+1))
			args = append(args, *filters.Department)
		}
		if filters.IsActive != nil {
			where = append(where, fmt.Sprintf("is_active = $%d", len(args)+1))
			args = append(args, *filters.IsActive)
		}
	}
	whereSQL := strings.Join(where, " AND ")

	countQ := `SELECT COUNT(*) FROM employees WHERE ` + whereSQL
	var total int64
	if err := r.sqlx.GetContext(ctx, &total, countQ, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	orderCol := "created_at"
	if filters != nil && filters.SortBy != "" {
		switch filters.SortBy {
		case "created_at", "department", "position", "hire_date", "employee_code":
			orderCol = filters.SortBy
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

	n := len(args)
	listQ := sqlSelectEmployeeBase + ` AND ` + whereSQL + fmt.Sprintf(` ORDER BY %s %s LIMIT $%d OFFSET $%d`, orderCol, dir, n+1, n+2)
	listArgs := append(append([]interface{}(nil), args...), limit, offset)

	var rows []EmployeeModel
	if err := r.sqlx.SelectContext(ctx, &rows, listQ, listArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}
	if err := r.enrichEmployeesWithUsers(ctx, rows); err != nil {
		return nil, 0, err
	}
	return r.employeesToDomain(rows), total, nil
}

func (r *PostgresEmployeeRepository) selectEmployeesSQLX(ctx context.Context, cond string, condArgs []interface{}, limit, offset int, errLabel string) ([]EmployeeModel, error) {
	q := sqlSelectEmployeeBase
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
	var rows []EmployeeModel
	if err := r.sqlx.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("%s: %w", errLabel, err)
	}
	return rows, nil
}

func (r *PostgresEmployeeRepository) enrichEmployeesWithUsers(ctx context.Context, emps []EmployeeModel) error {
	if r.sqlx == nil || len(emps) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, 0, len(emps))
	seen := make(map[uuid.UUID]struct{})
	for i := range emps {
		id := emps[i].UserID
		if id == uuid.Nil {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	m, err := fetchUserModelsByIDs(ctx, r.sqlx, ids)
	if err != nil {
		return fmt.Errorf("failed to load users: %w", err)
	}
	for i := range emps {
		if u, ok := m[emps[i].UserID]; ok {
			uCopy := u
			emps[i].User = &uCopy
		}
	}
	return nil
}

func (r *PostgresEmployeeRepository) employeesToDomain(rows []EmployeeModel) []*domain.Employee {
	out := make([]*domain.Employee, len(rows))
	for i := range rows {
		out[i] = r.toDomainEmployee(&rows[i])
	}
	return out
}

// toDomainEmployee converts database model to domain entity
func (r *PostgresEmployeeRepository) toDomainEmployee(dbEmployee *EmployeeModel) *domain.Employee {
	employee := &domain.Employee{
		ID:           dbEmployee.ID,
		UserID:       dbEmployee.UserID,
		EmployeeCode: dbEmployee.EmployeeCode,
		Position:     dbEmployee.Position,
		Department:   dbEmployee.Department,
		HireDate:     dbEmployee.HireDate,
		Salary:       dbEmployee.Salary,
		HoursPerWeek: dbEmployee.HoursPerWeek,
		IsActive:     dbEmployee.IsActive,
		CreatedAt:    dbEmployee.CreatedAt,
		UpdatedAt:    dbEmployee.UpdatedAt,
	}
	if dbEmployee.User != nil {
		employee.User = &domain.User{
			ID:        dbEmployee.User.ID,
			Email:     dbEmployee.User.Email,
			FirstName: dbEmployee.User.FirstName,
			LastName:  dbEmployee.User.LastName,
			Role:      dbEmployee.User.Role,
			IsActive:  dbEmployee.User.IsActive,
			CreatedAt: dbEmployee.User.CreatedAt,
			UpdatedAt: dbEmployee.User.UpdatedAt,
		}
	}
	return employee
}

// EmployeeModel represents the database table structure
type EmployeeModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" db:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" db:"user_id"`
	EmployeeCode string     `gorm:"uniqueIndex;not null" db:"employee_code"`
	Position     string     `gorm:"not null" db:"position"`
	Department   string     `gorm:"not null;index" db:"department"`
	HireDate     time.Time  `gorm:"not null" db:"hire_date"`
	Salary       float64    `gorm:"type:decimal(10,2)" db:"salary"`
	HoursPerWeek int        `gorm:"default:40" db:"hours_per_week"`
	IsActive     bool       `gorm:"default:true;index" db:"is_active"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" db:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" db:"deleted_at"`

	User *UserModel `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (EmployeeModel) TableName() string {
	return "employees"
}

// FindByID implements EmployeeRepository.FindByID
func (r *PostgresEmployeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	if r.sqlx != nil {
		var row EmployeeModel
		err := r.sqlx.GetContext(ctx, &row, sqlSelectEmployeeBase+` AND id = $1`, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrEmployeeNotFound
			}
			return nil, fmt.Errorf("failed to get employee by ID: %w", err)
		}
		rows := []EmployeeModel{row}
		if err := r.enrichEmployeesWithUsers(ctx, rows); err != nil {
			return nil, err
		}
		return r.toDomainEmployee(&rows[0]), nil
	}
	var dbEmployee EmployeeModel
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("employees.id = ? AND employees.deleted_at IS NULL", id).
		First(&dbEmployee).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee by ID: %w", err)
	}
	return r.toDomainEmployee(&dbEmployee), nil
}
