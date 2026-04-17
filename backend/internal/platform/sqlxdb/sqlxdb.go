// Package sqlxdb opens PostgreSQL connections with sqlx + lib/pq (template Phase 2).
// Use this for new code paths; GORM remains until vertical migrations complete.
package sqlxdb

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // driver name "postgres"
)

// Open returns a sqlx DB connected with the postgres driver. dsn is a lib/pq connection string
// (same shape as DATABASE_URL used by GORM, e.g. postgres://user:pass@host:5432/db?sslmode=disable).
func Open(dsn string) (*sqlx.DB, error) {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return nil, fmt.Errorf("sqlxdb: empty dsn")
	}
	return sqlx.Connect("postgres", dsn)
}
