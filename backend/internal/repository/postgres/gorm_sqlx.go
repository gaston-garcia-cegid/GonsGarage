package postgres

import (
	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/platform/sqlxdb"
)

// sqlxFromGORM returns a sqlx handle over the same pool as GORM when the dialector is PostgreSQL.
func sqlxFromGORM(db *gorm.DB) *sqlx.DB {
	sqlDB, err := db.DB()
	if err != nil || sqlDB == nil {
		return nil
	}
	switch db.Dialector.Name() {
	case "postgres", "pgx":
		return sqlxdb.WrapPostgres(sqlDB)
	default:
		return nil
	}
}
