// Command seed-test-client creates a fixed demo client user if it does not exist (idempotent).
// Usage (from backend/): go run ./cmd/seed-test-client
//
// Env: DATABASE_URL (optional), SEED_CLIENT_EMAIL, SEED_CLIENT_PASSWORD
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	postgresRepo "github.com/gaston-garcia-cegid/gonsgarage/internal/repository/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://admindb:gonsgarage123@localhost:5432/gonsgarage?sslmode=disable"
	}

	email := os.Getenv("SEED_CLIENT_EMAIL")
	if email == "" {
		email = "cliente.demo@gonsgarage.local"
	}
	password := os.Getenv("SEED_CLIENT_PASSWORD")
	if password == "" {
		password = "ClienteDemo123"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("conexión a base de datos: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("migrate users: %v", err)
	}

	repo := postgresRepo.NewPostgresUserRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	existing, err := repo.GetByEmail(ctx, email)
	if err != nil {
		log.Fatalf("consultar email: %v", err)
	}
	if existing != nil {
		log.Printf("El usuario %s ya existe (rol=%s). No se creó nada.", email, existing.Role)
		os.Exit(0)
	}

	user, err := domain.NewUser(email, password, "Cliente", "Demo", domain.RoleClient)
	if err != nil {
		log.Fatalf("crear modelo de usuario: %v", err)
	}

	if err := repo.Create(ctx, user); err != nil {
		log.Fatalf("insertar usuario: %v", err)
	}

	log.Printf("Usuario cliente de prueba creado correctamente.")
	log.Printf("  Email:    %s", email)
	log.Printf("  Password: %s  (definir SEED_CLIENT_PASSWORD para otra)", password)
	log.Printf("Iniciá sesión en http://localhost:3000/auth/login con esas credenciales.")
}
