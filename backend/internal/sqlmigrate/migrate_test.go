package sqlmigrate

import (
	"strings"
	"testing"
)

func TestPostgresURLWithMigrationsTable_appendsDefault(t *testing.T) {
	t.Parallel()
	in := "postgres://u:p@localhost:5432/db?sslmode=disable"
	out, err := postgresURLWithMigrationsTable(in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "x-migrations-table=gon_golang_migrate") {
		t.Fatalf("expected x-migrations-table in url, got %q", out)
	}
}

func TestPostgresURLWithMigrationsTable_preservesExplicit(t *testing.T) {
	t.Parallel()
	in := "postgres://u:p@localhost:5432/db?sslmode=disable&x-migrations-table=my_custom"
	out, err := postgresURLWithMigrationsTable(in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "x-migrations-table=my_custom") {
		t.Fatalf("got %q", out)
	}
	if strings.Count(out, "x-migrations-table") > 1 {
		t.Fatalf("duplicate param: %q", out)
	}
}
