// Package sqlxdb integrates sqlx with the existing PostgreSQL pool (template Phase 2).
// GORM owns the *sql.DB; we wrap it for sqlx helpers and readiness without a second pool.
package sqlxdb

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// WrapPostgres returns a sqlx handle over db using the postgres driver name (pgx/libpq via GORM).
func WrapPostgres(db *sql.DB) *sqlx.DB {
	if db == nil {
		return nil
	}
	return sqlx.NewDb(db, "postgres")
}
