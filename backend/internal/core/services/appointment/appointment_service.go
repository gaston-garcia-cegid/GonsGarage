package appointment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
)

type AppointmentService struct {
	repo      ports.AppointmentRepository
	userRepo  ports.UserRepository
	cacheRepo ports.CacheRepository
}

func NewAppointmentService(
	repo ports.AppointmentRepository,
	userRepo ports.UserRepository,
	cacheRepo ports.CacheRepository) *AppointmentService {
	return &AppointmentService{
		repo:      repo,
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
	}
}

// CreateAppointment creates a new appointment
func (s *AppointmentService) CreateAppointment(
	ctx context.Context,
	appointment *domain.Appointment,
	requestingUserID uuid.UUID) (*domain.Appointment, error) {

	queryCtx, cancel := context.WithTimeout(ctx, 10*time.Second) // ✅ Increase timeout
	defer cancel()

	requestingUser, err := s.userRepo.GetByID(queryCtx, requestingUserID)
	if err != nil {
		log.Printf("failed to get requesting user: userID=%s, error=%v", requestingUserID, err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if requestingUser == nil {
		log.Printf("user not found: userID=%s", requestingUserID)
		return nil, domain.ErrUserNotFound
	}

	// ✅ Set metadata
	appointment.ID = uuid.New()
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()

	// Get the requesting user to check permissions
	// requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	// if err != nil {
	// 	return nil, err
	// }

	// Check if the user has permission to create appointments
	// if !requestingUser.HasPermission(domain.PermissionCreateAppointment) {
	// 	return nil, domain.ErrPermissionDenied
	// }

	// Create the appointment
	if err := s.repo.Create(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

// UpdateAppointment updates an existing appointment
func (s *AppointmentService) UpdateAppointment(ctx context.Context, appointment *domain.Appointment, requestingUserID uuid.UUID) (*domain.Appointment, error) {
	// Get the requesting user to check permissions
	// requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	// if err != nil {
	// 	return nil, err
	// }

	// // Check if the user has permission to update appointments
	// if !requestingUser.HasPermission(domain.PermissionUpdateAppointment) {
	// 	return nil, domain.ErrPermissionDenied
	// }

	// Update the appointment
	if err := s.repo.Update(ctx, appointment); err != nil {
		return nil, err
	}

	return appointment, nil
}

// CancelAppointment cancels an existing appointment
// func (s *AppointmentService) CancelAppointment(ctx context.Context, appointmentID uuid.UUID, requestingUserID uuid.UUID) error {
// 	// Get the requesting user to check permissions
// 	// requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// // Check if the user has permission to cancel appointments
// 	// if !requestingUser.HasPermission(domain.PermissionCancelAppointment) {
// 	// 	return domain.ErrPermissionDenied
// 	// }

// 	// Cancel the appointment
// 	err = s.repo.Cancel(ctx, appointmentID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// GetAppointment retrieves an appointment by ID
func (s *AppointmentService) GetAppointment(ctx context.Context, appointmentID uuid.UUID, requestingUserID uuid.UUID) (*domain.Appointment, error) {
	// Get the requesting user to check permissions
	// requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	// if err != nil {
	// 	return nil, err
	// }

	// // Check if the user has permission to view appointments
	// if !requestingUser.HasPermission(domain.PermissionViewAppointment) {
	// 	return nil, domain.ErrPermissionDenied
	// }

	// Get the appointment
	appointment, err := s.repo.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}

	return appointment, nil
}

// ListAppointments lists appointments with optional filters
func (s *AppointmentService) ListAppointments(ctx context.Context, filters *ports.AppointmentFilters) ([]*domain.Appointment, int64, error) {
	// Get the requesting user to check permissions
	// requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	// if err != nil {
	// 	return nil, 0, err
	// }

	// // Check if the user has permission to view appointments
	// if !requestingUser.HasPermission(domain.PermissionViewAppointment) {
	// 	return nil, 0, domain.ErrPermissionDenied
	// }

	if filters == nil {
		filters = &ports.AppointmentFilters{
			Limit:     10,
			Offset:    0,
			SortBy:    "created_at",
			SortOrder: "DESC",
		}
	}

	if filters.Limit == 0 {
		filters.Limit = 10
	}

	if filters.SortBy == "" {
		filters.SortBy = "created_at"
	}

	if filters.SortOrder == "" {
		filters.SortOrder = "DESC"
	}

	// List the appointments
	appointments, total, err := s.repo.List(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

// DeleteAppointment deletes an appointment by ID
func (s *AppointmentService) DeleteAppointment(ctx context.Context, appointmentID uuid.UUID, requestingUserID uuid.UUID) error {
	// Get the requesting user to check permissions
	// requestingUser, err := s.userRepo.GetByID(ctx, requestingUserID)
	// if err != nil {
	// 	return err
	// }

	// // Check if the user has permission to delete appointments
	// if !requestingUser.HasPermission(domain.PermissionDeleteAppointment) {
	// 	return domain.ErrPermissionDenied
	// }

	// Delete the appointment
	if err := s.repo.Delete(ctx, appointmentID); err != nil {
		return err
	}

	return nil
}
