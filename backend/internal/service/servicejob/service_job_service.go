package servicejob

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/google/uuid"
)

// Service implements workshop service-job use cases.
type Service struct {
	jobRepo  ports.ServiceJobRepository
	carRepo  ports.CarRepository
	userRepo ports.UserRepository
}

func NewService(jobRepo ports.ServiceJobRepository, carRepo ports.CarRepository, userRepo ports.UserRepository) *Service {
	return &Service{jobRepo: jobRepo, carRepo: carRepo, userRepo: userRepo}
}

func (s *Service) requireWorkshopUser(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if u == nil || !u.IsEmployee() {
		return nil, domain.ErrUnauthorizedAccess
	}
	return u, nil
}

func (s *Service) canAccessCar(ctx context.Context, user *domain.User, carID uuid.UUID) (*domain.Car, error) {
	car, err := s.carRepo.GetByID(ctx, carID)
	if err != nil {
		return nil, fmt.Errorf("car: %w", err)
	}
	if car == nil {
		return nil, domain.ErrCarNotFound
	}
	if user.IsClient() && car.OwnerID != user.ID {
		return nil, domain.ErrUnauthorizedAccess
	}
	return car, nil
}

// CreateServiceJob starts a new visit; only workshop staff.
func (s *Service) CreateServiceJob(ctx context.Context, carID uuid.UUID, userID uuid.UUID) (*domain.ServiceJob, error) {
	u, err := s.requireWorkshopUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if _, err := s.canAccessCar(ctx, u, carID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	job := &domain.ServiceJob{
		ID:             uuid.New(),
		CarID:          carID,
		Status:         domain.ServiceJobStatusOpen,
		OpenedByUserID: userID,
		OpenedAt:       now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

// GetWithDetails returns job, reception, handover if present.
func (s *Service) GetWithDetails(ctx context.Context, jobID uuid.UUID, userID uuid.UUID) (*domain.ServiceJob, *domain.ServiceJobReception, *domain.ServiceJobHandover, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get user: %w", err)
	}
	if u == nil {
		return nil, nil, nil, domain.ErrUnauthorizedAccess
	}
	j, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, nil, nil, err
	}
	if _, err := s.canAccessCar(ctx, u, j.CarID); err != nil {
		return nil, nil, nil, err
	}
	rec, err := s.jobRepo.GetReception(ctx, jobID)
	if err != nil {
		return nil, nil, nil, err
	}
	ho, err := s.jobRepo.GetHandover(ctx, jobID)
	if err != nil {
		return nil, nil, nil, err
	}
	return j, rec, ho, nil
}

// ListByCarID returns visits for a car (client: own car; staff: any in catalog).
func (s *Service) ListByCarID(ctx context.Context, carID uuid.UUID, userID uuid.UUID) ([]*domain.ServiceJob, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if u == nil {
		return nil, domain.ErrUnauthorizedAccess
	}
	if _, err := s.canAccessCar(ctx, u, carID); err != nil {
		return nil, err
	}
	return s.jobRepo.ListByCarID(ctx, carID)
}

// SaveReceptionInput is validated service-layer input.
type SaveReceptionInput struct {
	OdometerKM   int
	OilLevel     string
	CoolantLevel string
	TiresNote    string
	GeneralNotes string
}

func (s *Service) SaveReception(ctx context.Context, jobID uuid.UUID, in SaveReceptionInput, userID uuid.UUID) (*domain.ServiceJobReception, error) {
	u, err := s.requireWorkshopUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if in.OdometerKM < 0 {
		return nil, domain.ErrInvalidServiceJobData
	}
	j, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if j.Status == domain.ServiceJobStatusClosed || j.Status == domain.ServiceJobStatusCancelled {
		return nil, domain.ErrInvalidServiceJobData
	}
	if _, err := s.canAccessCar(ctx, u, j.CarID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	r := &domain.ServiceJobReception{
		ServiceJobID:     jobID,
		OdometerKM:       in.OdometerKM,
		OilLevel:         in.OilLevel,
		CoolantLevel:     in.CoolantLevel,
		TiresNote:        in.TiresNote,
		GeneralNotes:     in.GeneralNotes,
		RecordedByUserID: userID,
		RecordedAt:       now,
		SchemaVersion:    1,
	}
	if err := s.jobRepo.SaveReception(ctx, r); err != nil {
		return nil, err
	}
	j.Status = domain.ServiceJobStatusInProgress
	j.UpdatedAt = now
	_ = s.jobRepo.Update(ctx, j) // best-effort status
	return r, nil
}

// SaveHandoverInput is service-layer handover DTO.
type SaveHandoverInput struct {
	OdometerKM   int
	TiresNote    string
	GeneralNotes string
}

func (s *Service) SaveHandover(ctx context.Context, jobID uuid.UUID, in SaveHandoverInput, userID uuid.UUID) (*domain.ServiceJobHandover, error) {
	u, err := s.requireWorkshopUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if in.OdometerKM < 0 {
		return nil, domain.ErrInvalidServiceJobData
	}
	j, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if j.Status == domain.ServiceJobStatusClosed {
		existing, _ := s.jobRepo.GetHandover(ctx, jobID)
		if existing != nil {
			return existing, nil
		}
	}
	if j.Status == domain.ServiceJobStatusCancelled {
		return nil, domain.ErrInvalidServiceJobData
	}
	if _, err := s.canAccessCar(ctx, u, j.CarID); err != nil {
		return nil, err
	}
	prev, err := s.jobRepo.GetReception(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if prev == nil {
		return nil, domain.ErrReceptionRequiredBeforeHandover
	}
	now := time.Now().UTC()
	h := &domain.ServiceJobHandover{
		ServiceJobID:     jobID,
		OdometerKM:       in.OdometerKM,
		TiresNote:        in.TiresNote,
		GeneralNotes:     in.GeneralNotes,
		RecordedByUserID: userID,
		RecordedAt:       now,
		SchemaVersion:    1,
	}
	if err := s.jobRepo.SaveHandover(ctx, h); err != nil {
		return nil, err
	}
	j.Status = domain.ServiceJobStatusClosed
	j.ClosedAt = &now
	j.UpdatedAt = now
	if err := s.jobRepo.Update(ctx, j); err != nil {
		return nil, err
	}
	return h, nil
}
