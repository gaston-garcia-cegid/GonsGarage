package sqlmigrate

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// defaultMigrationsTable evita colisión con otras herramientas que crean
// public.schema_migrations sin la columna dirty (p. ej. Rails u otros migradores).
const defaultMigrationsTable = "gon_golang_migrate"

func postgresURLWithMigrationsTable(databaseURL string) (string, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", fmt.Errorf("parse database url: %w", err)
	}
	q := u.Query()
	if q.Get("x-migrations-table") == "" {
		q.Set("x-migrations-table", defaultMigrationsTable)
		u.RawQuery = q.Encode()
	}
	return u.String(), nil
}

// tryFixDirtyIfEnabled limpia el flag dirty de golang-migrate (p. ej. tras un fallo por BOM en SQL).
// Solo si MIGRATE_AUTO_FIX_DIRTY=true y GIN_MODE≠release. No usar en producción sin revisar el estado real de la BD.
func tryFixDirtyIfEnabled(m *migrate.Migrate) error {
	if !strings.EqualFold(strings.TrimSpace(os.Getenv("MIGRATE_AUTO_FIX_DIRTY")), "true") {
		return nil
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("GIN_MODE")), "release") {
		return fmt.Errorf("MIGRATE_AUTO_FIX_DIRTY no está permitido con GIN_MODE=release")
	}
	v, dirty, verr := m.Version()
	if verr != nil {
		if errors.Is(verr, migrate.ErrNilVersion) {
			return nil
		}
		return fmt.Errorf("migrate version check: %w", verr)
	}
	if !dirty {
		return nil
	}
	if err := m.Force(int(v)); err != nil {
		return fmt.Errorf("migrate force (clear dirty at version %d): %w", v, err)
	}
	return nil
}

func upWithHint(err error) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	if strings.Contains(s, "Dirty database") {
		return fmt.Errorf("%w\nhint: en dev una vez MIGRATE_AUTO_FIX_DIRTY=true (sin GIN_MODE=release), o SQL: UPDATE %s SET dirty=false WHERE dirty=true",
			err, defaultMigrationsTable)
	}
	return err
}

// Up applies all pending migrations from migrationsDir (filesystem path).
// databaseURL is a postgres DSN accepted by golang-migrate (e.g. postgres://user:pass@host:5432/db?sslmode=disable).
// Si el DSN no incluye x-migrations-table, se usa gon_golang_migrate para no chocar con schema_migrations de otros stacks.
func Up(databaseURL, migrationsDir string) error {
	dbURL, err := postgresURLWithMigrationsTable(databaseURL)
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("migrations path: %w", err)
	}
	src := "file://" + filepath.ToSlash(abs)
	m, err := migrate.New(src, dbURL)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	if err := tryFixDirtyIfEnabled(m); err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return upWithHint(err)
	}
	return nil
}
