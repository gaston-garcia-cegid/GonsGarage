package repositories

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

type CarRepository interface {
	Create(ctx context.Context, car *domain.Car) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error)
	GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Car, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Car, error)
	Update(ctx context.Context, car *domain.Car) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error)
	GetDeletedByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error)
	Restore(ctx context.Context, id uuid.UUID) error
}

type RepairRepository interface {
	Create(ctx context.Context, repair *domain.Repair) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error)
	GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Repair, error)
	Update(ctx context.Context, repair *domain.Repair) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *domain.Appointment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error)
	GetByClientID(ctx context.Context, clientID uuid.UUID) ([]*domain.Appointment, error)
	GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Appointment, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Appointment, error)
	Update(ctx context.Context, appointment *domain.Appointment) error
	Delete(ctx context.Context, id uuid.UUID) error
}
