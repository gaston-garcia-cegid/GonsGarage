// Package repository is the persistence entrypoint (template §1).
// Phase 1: aliases to internal/adapters/repository/postgres constructors.
package repository

import postgres "github.com/gaston-garcia-cegid/gonsgarage/internal/adapters/repository/postgres"

var (
	NewPostgresUserRepository        = postgres.NewPostgresUserRepository
	NewPostgresEmployeeRepository    = postgres.NewPostgresEmployeeRepository
	NewPostgresCarRepository         = postgres.NewPostgresCarRepository
	NewPostgresRepairRepository      = postgres.NewPostgresRepairRepository
	NewPostgresAppointmentRepository = postgres.NewPostgresAppointmentRepository
	NewWorkshopRepository            = postgres.NewWorkshopRepository
	NewAccountingRepository          = postgres.NewAccountingRepository
	NewNotificationRepository        = postgres.NewNotificationRepository
)
