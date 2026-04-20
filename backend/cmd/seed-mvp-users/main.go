// Command seed-mvp-users creates idempotent demo users for admin, manager, and employee (JWT roles).
// Client demo remains: go run ./cmd/seed-test-client
//
// Usage (from backend/): go run ./cmd/seed-mvp-users
//
// Env: DATABASE_URL (optional). Per role: SEED_ADMIN_EMAIL|PASSWORD, SEED_MANAGER_EMAIL|PASSWORD, SEED_EMPLOYEE_EMAIL|PASSWORD (each optional; defaults *.gonsgarage.local).
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	postgresRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/repository/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type seedDef struct {
	role          string
	first, last   string
	envEmail      string
	defaultEmail  string
	envPassword   string
	defaultPass   string
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admindb:gonsgarage123@localhost:5432/gonsgarage?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	repo := postgresRepo.NewPostgresUserRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	defs := []seedDef{
		{
			role:         domain.RoleAdmin,
			first:        "Admin",
			last:         "Demo",
			envEmail:     "SEED_ADMIN_EMAIL",
			defaultEmail: "admin.demo@gonsgarage.local",
			envPassword:  "SEED_ADMIN_PASSWORD",
			defaultPass:  "AdminDemo123",
		},
		{
			role:         domain.RoleManager,
			first:        "Gestor",
			last:         "Demo",
			envEmail:     "SEED_MANAGER_EMAIL",
			defaultEmail: "manager.demo@gonsgarage.local",
			envPassword:  "SEED_MANAGER_PASSWORD",
			defaultPass:  "ManagerDemo123",
		},
		{
			role:         domain.RoleEmployee,
			first:        "Mecânico",
			last:         "Demo",
			envEmail:     "SEED_EMPLOYEE_EMAIL",
			defaultEmail: "employee.demo@gonsgarage.local",
			envPassword:  "SEED_EMPLOYEE_PASSWORD",
			defaultPass:  "EmployeeDemo123",
		},
	}

	for _, d := range defs {
		email := os.Getenv(d.envEmail)
		if email == "" {
			email = d.defaultEmail
		}
		password := os.Getenv(d.envPassword)
		if password == "" {
			password = d.defaultPass
		}

		existing, err := repo.GetByEmail(ctx, email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			log.Fatalf("%s (%s): %v", d.role, email, err)
		}
		if existing != nil {
			log.Printf("%s: user %s already exists (role=%s). Skip.", d.role, email, existing.Role)
			continue
		}

		user, err := domain.NewUser(email, password, d.first, d.last, d.role)
		if err != nil {
			log.Fatalf("%s: NewUser: %v", d.role, err)
		}
		if err := repo.Create(ctx, user); err != nil {
			log.Fatalf("%s: Create: %v", d.role, err)
		}
		log.Printf("%s: created %s", d.role, email)
	}

	log.Println("seed-mvp-users done. Client demo: go run ./cmd/seed-test-client")
}
