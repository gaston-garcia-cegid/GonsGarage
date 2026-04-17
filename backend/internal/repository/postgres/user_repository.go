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

const (
	sqlSelectUserBase = `SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at, deleted_at
FROM users WHERE deleted_at IS NULL`
)

// PostgresUserRepository implements UserRepository interface using PostgreSQL
type PostgresUserRepository struct {
	db   *gorm.DB
	sqlx *sqlx.DB // same pool as GORM (Phase 2: hot paths with sqlx)
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *gorm.DB) ports.UserRepository {
	return &PostgresUserRepository{db: db, sqlx: sqlxFromGORM(db)}
}

// Create implements UserRepository.Create
func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	if r.sqlx != nil {
		return r.createSQLX(ctx, user)
	}
	dbUser := &UserModel{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.Password,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		IsActive:     user.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

func (r *PostgresUserRepository) createSQLX(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC()
	const q = `INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING created_at, updated_at`
	row := r.sqlx.QueryRowxContext(ctx, q,
		user.ID, user.Email, user.Password, user.FirstName, user.LastName, user.Role, user.IsActive, now, now,
	)
	if err := row.Scan(&user.CreatedAt, &user.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID implements UserRepository.GetByID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if r.sqlx != nil {
		return r.getByIDSQLX(ctx, id)
	}
	var dbUser UserModel
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := r.db.WithContext(queryCtx).Where("id = ? AND deleted_at IS NULL", id).First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		if err == context.Canceled || err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database query timeout: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return r.toDomainUser(&dbUser), nil
}

func (r *PostgresUserRepository) getByIDSQLX(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var row UserModel
	q := sqlSelectUserBase + ` AND id = $1`
	err := r.sqlx.GetContext(queryCtx, &row, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timeout: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return r.toDomainUser(&row), nil
}

// GetByEmail implements UserRepository.GetByEmail
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if r.sqlx != nil {
		var row UserModel
		q := sqlSelectUserBase + ` AND email = $1`
		err := r.sqlx.GetContext(ctx, &row, q, email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, domain.ErrUserNotFound
			}
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
		return r.toDomainUser(&row), nil
	}
	var dbUser UserModel
	err := r.db.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return r.toDomainUser(&dbUser), nil
}

// Update implements UserRepository.Update
func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	if r.sqlx != nil {
		return r.updateSQLX(ctx, user)
	}
	dbUser := &UserModel{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.Password,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		IsActive:     user.IsActive,
		UpdatedAt:    time.Now(),
	}
	result := r.db.WithContext(ctx).Model(dbUser).Where("id = ?", user.ID).Updates(dbUser)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	user.UpdatedAt = dbUser.UpdatedAt
	return nil
}

func (r *PostgresUserRepository) updateSQLX(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC()
	const q = `UPDATE users SET
email = $1, password_hash = $2, first_name = $3, last_name = $4, role = $5, is_active = $6, updated_at = $7
WHERE id = $8 AND deleted_at IS NULL`
	res, err := r.sqlx.ExecContext(ctx, q,
		user.Email, user.Password, user.FirstName, user.LastName, user.Role, user.IsActive, now, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read rows affected: %w", err)
	}
	if n == 0 {
		return domain.ErrUserNotFound
	}
	user.UpdatedAt = now
	return nil
}

// Delete implements UserRepository.Delete (soft delete)
func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if r.sqlx != nil {
		const q = `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
		res, err := r.sqlx.ExecContext(ctx, q, time.Now().UTC(), id)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read rows affected: %w", err)
		}
		if n == 0 {
			return domain.ErrUserNotFound
		}
		return nil
	}
	result := r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// List implements UserRepository.List
func (r *PostgresUserRepository) List(ctx context.Context, limit int, offset int) ([]*domain.User, error) {
	if r.sqlx != nil {
		return r.selectUsersSQLX(ctx, "", nil, limit, offset, "failed to list users")
	}
	var dbUsers []UserModel
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&dbUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toDomainUser(&dbUser)
	}
	return users, nil
}

// selectUsersSQLX runs sqlSelectUserBase + optional AND cond with placeholder args, then ORDER BY created_at DESC.
func (r *PostgresUserRepository) selectUsersSQLX(ctx context.Context, cond string, condArgs []interface{}, limit, offset int, errLabel string) ([]*domain.User, error) {
	q := sqlSelectUserBase
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
	var rows []UserModel
	if err := r.sqlx.SelectContext(ctx, &rows, q, args...); err != nil {
		return nil, fmt.Errorf("%s: %w", errLabel, err)
	}
	out := make([]*domain.User, len(rows))
	for i := range rows {
		out[i] = r.toDomainUser(&rows[i])
	}
	return out, nil
}

// EmailExists implements UserRepository.EmailExists
func (r *PostgresUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	if r.sqlx != nil {
		const q = `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND deleted_at IS NULL)`
		var exists bool
		if err := r.sqlx.GetContext(ctx, &exists, q, email); err != nil {
			return false, fmt.Errorf("failed to check email existence: %w", err)
		}
		return exists, nil
	}
	var count int64
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ? AND deleted_at IS NULL", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

// GetActiveUsers implements UserRepository.GetActiveUsers
func (r *PostgresUserRepository) GetActiveUsers(ctx context.Context, limit int, offset int) ([]*domain.User, error) {
	if r.sqlx != nil {
		return r.selectUsersSQLX(ctx, "is_active = $1", []interface{}{true}, limit, offset, "failed to list active users")
	}
	var dbUsers []UserModel
	query := r.db.WithContext(ctx).Where("is_active = ? AND deleted_at IS NULL", true).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&dbUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list active users: %w", err)
	}
	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toDomainUser(&dbUser)
	}
	return users, nil
}

// GetByRole implements UserRepository.GetByRole
func (r *PostgresUserRepository) GetByRole(ctx context.Context, role string, limit int, offset int) ([]*domain.User, error) {
	if r.sqlx != nil {
		return r.selectUsersSQLX(ctx, "role = $1", []interface{}{role}, limit, offset, "failed to get users by role")
	}
	var dbUsers []UserModel
	query := r.db.WithContext(ctx).Where("role = ? AND deleted_at IS NULL", role).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&dbUsers).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	users := make([]*domain.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toDomainUser(&dbUser)
	}
	return users, nil
}

// UpdatePassword implements UserRepository.UpdatePassword
func (r *PostgresUserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error {
	if r.sqlx != nil {
		const q = `UPDATE users SET password_hash = $1, updated_at = $2
WHERE id = $3 AND deleted_at IS NULL`
		res, err := r.sqlx.ExecContext(ctx, q, newPasswordHash, time.Now().UTC(), userID)
		if err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to read rows affected: %w", err)
		}
		if n == 0 {
			return domain.ErrUserNotFound
		}
		return nil
	}
	result := r.db.WithContext(ctx).Model(&UserModel{}).
		Where("id = ? AND deleted_at IS NULL", userID).
		Update("password_hash", newPasswordHash)
	if result.Error != nil {
		return fmt.Errorf("failed to update password: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

// toDomainUser converts database model to domain entity
func (r *PostgresUserRepository) toDomainUser(dbUser *UserModel) *domain.User {
	if dbUser == nil {
		return nil
	}

	return &domain.User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Password:  dbUser.PasswordHash,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Role:      dbUser.Role,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

// UserModel represents the database table structure
type UserModel struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" db:"id"`
	Email        string     `gorm:"uniqueIndex;not null" db:"email"`
	PasswordHash string     `gorm:"not null" db:"password_hash"`
	FirstName    string     `gorm:"not null" db:"first_name"`
	LastName     string     `gorm:"not null" db:"last_name"`
	Role         string     `gorm:"not null;default:'employee'" db:"role"`
	IsActive     bool       `gorm:"default:true" db:"is_active"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" db:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" db:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" db:"deleted_at"`
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}
