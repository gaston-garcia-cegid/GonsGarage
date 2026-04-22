package servicejob

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tUser map[uuid.UUID]*domain.User

func (m tUser) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := m[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}
func (m tUser) Create(context.Context, *domain.User) error                     { return nil }
func (m tUser) GetByEmail(context.Context, string) (*domain.User, error)       { return nil, nil }
func (m tUser) GetByRole(context.Context, string, int, int) ([]*domain.User, error) {
	return nil, nil
}
func (m tUser) List(context.Context, int, int) ([]*domain.User, error)   { return nil, nil }
func (m tUser) Update(context.Context, *domain.User) error                 { return nil }
func (m tUser) Delete(context.Context, uuid.UUID) error                 { return nil }
func (m tUser) UpdatePassword(context.Context, uuid.UUID, string) error  { return nil }
func (m tUser) GetActiveUsers(context.Context, int, int) ([]*domain.User, error) { return nil, nil }

type tCar map[uuid.UUID]*domain.Car

func (c tCar) GetByID(_ context.Context, id uuid.UUID) (*domain.Car, error) {
	x, ok := c[id]
	if !ok {
		return nil, domain.ErrCarNotFound
	}
	return x, nil
}
func (c tCar) Create(context.Context, *domain.Car) error                     { return nil }
func (c tCar) GetByOwnerID(context.Context, uuid.UUID) ([]*domain.Car, error) { return nil, nil }
func (c tCar) GetByLicensePlate(context.Context, string) (*domain.Car, error) {
	return nil, nil
}
func (c tCar) List(context.Context, int, int) ([]*domain.Car, error)  { return nil, nil }
func (c tCar) Update(context.Context, *domain.Car) error               { return nil }
func (c tCar) Delete(context.Context, uuid.UUID) error                 { return nil }
func (c tCar) GetWithRepairs(context.Context, uuid.UUID) (*domain.Car, error) { return nil, nil }
func (c tCar) GetDeletedByLicensePlate(context.Context, string) (*domain.Car, error) {
	return nil, nil
}
func (c tCar) Restore(context.Context, uuid.UUID) error { return nil }

type stubJobRepo struct {
	created  []*domain.ServiceJob
	byID     map[uuid.UUID]*domain.ServiceJob
	byCar    map[uuid.UUID][]*domain.ServiceJob
	rec      map[uuid.UUID]*domain.ServiceJobReception
	handover map[uuid.UUID]*domain.ServiceJobHandover
}

func (s *stubJobRepo) Create(_ context.Context, j *domain.ServiceJob) error {
	if s.byID == nil {
		s.byID = make(map[uuid.UUID]*domain.ServiceJob)
	}
	s.byID[j.ID] = j
	s.created = append(s.created, j)
	if s.byCar != nil {
		s.byCar[j.CarID] = append(s.byCar[j.CarID], j)
	}
	return nil
}

func (s *stubJobRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.ServiceJob, error) {
	j, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrServiceJobNotFound
	}
	return j, nil
}

func (s *stubJobRepo) Update(_ context.Context, j *domain.ServiceJob) error {
	if s.byID == nil {
		return domain.ErrServiceJobNotFound
	}
	if _, ok := s.byID[j.ID]; !ok {
		return domain.ErrServiceJobNotFound
	}
	s.byID[j.ID] = j
	return nil
}

func (s *stubJobRepo) ListByCarID(_ context.Context, carID uuid.UUID) ([]*domain.ServiceJob, error) {
	if s.byCar == nil {
		return nil, nil
	}
	return s.byCar[carID], nil
}

func (s *stubJobRepo) ListByOpenedOn(_ context.Context, day time.Time) ([]*domain.ServiceJob, error) {
	y, m, d := day.UTC().Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	var out []*domain.ServiceJob
	for _, j := range s.byID {
		if j.OpenedAt.Before(start) || !j.OpenedAt.Before(end) {
			continue
		}
		cj := *j
		out = append(out, &cj)
	}
	if out == nil {
		out = []*domain.ServiceJob{}
	}
	return out, nil
}

func (s *stubJobRepo) SaveReception(_ context.Context, r *domain.ServiceJobReception) error {
	if s.rec == nil {
		s.rec = make(map[uuid.UUID]*domain.ServiceJobReception)
	}
	s.rec[r.ServiceJobID] = r
	return nil
}

func (s *stubJobRepo) GetReception(_ context.Context, id uuid.UUID) (*domain.ServiceJobReception, error) {
	if s.rec == nil {
		return nil, nil
	}
	return s.rec[id], nil
}

func (s *stubJobRepo) SaveHandover(_ context.Context, h *domain.ServiceJobHandover) error {
	if s.handover == nil {
		s.handover = make(map[uuid.UUID]*domain.ServiceJobHandover)
	}
	s.handover[h.ServiceJobID] = h
	return nil
}

func (s *stubJobRepo) GetHandover(_ context.Context, id uuid.UUID) (*domain.ServiceJobHandover, error) {
	if s.handover == nil {
		return nil, nil
	}
	return s.handover[id], nil
}

func TestService_CreateServiceJob_ClientDenied(t *testing.T) {
	t.Parallel()
	carID := uuid.New()
	clientID := uuid.New()
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{}, byCar: map[uuid.UUID][]*domain.ServiceJob{carID: {}}}
	cu, _ := domain.NewUser("c@t", "p", "C", "L", domain.RoleClient)
	cu.ID = clientID
	s := NewService(&je, tCar{carID: {ID: carID, OwnerID: clientID}}, tUser{clientID: cu}, nil)

	out, err := s.CreateServiceJob(context.Background(), carID, clientID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestService_CreateServiceJob_EmployeeOK(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	carID := uuid.New()
	ownerID := uuid.New()
	emp, _ := domain.NewUser("e@t", "p", "E", "E", domain.RoleEmployee)
	emp.ID = empID
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{}, byCar: map[uuid.UUID][]*domain.ServiceJob{carID: {}}}
	s := NewService(&je, tCar{carID: {ID: carID, OwnerID: ownerID}}, tUser{empID: emp}, nil)

	out, err := s.CreateServiceJob(context.Background(), carID, empID)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, domain.ServiceJobStatusOpen, out.Status)
	assert.Equal(t, carID, out.CarID)
	require.Len(t, je.created, 1)
}

func TestService_SaveHandover_WithoutReception(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	jobID := uuid.New()
	carID := uuid.New()
	emp, _ := domain.NewUser("e@t", "p", "E", "E", domain.RoleEmployee)
	emp.ID = empID
	j := &domain.ServiceJob{ID: jobID, CarID: carID, Status: domain.ServiceJobStatusInProgress, OpenedByUserID: empID, OpenedAt: time.Now().UTC()}
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{jobID: j}, byCar: map[uuid.UUID][]*domain.ServiceJob{carID: {j}}, rec: map[uuid.UUID]*domain.ServiceJobReception{}}
	s := NewService(&je, tCar{carID: {ID: carID, OwnerID: uuid.New()}}, tUser{empID: emp}, nil)

	_, err := s.SaveHandover(context.Background(), jobID, SaveHandoverInput{OdometerKM: 10}, empID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrReceptionRequiredBeforeHandover)
}

func TestService_SaveReception_Then_Handover_Closed(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	jobID := uuid.New()
	carID := uuid.New()
	emp, _ := domain.NewUser("e@t", "p", "E", "E", domain.RoleEmployee)
	emp.ID = empID
	j := &domain.ServiceJob{ID: jobID, CarID: carID, Status: domain.ServiceJobStatusOpen, OpenedByUserID: empID, OpenedAt: time.Now().UTC(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{jobID: j}, byCar: map[uuid.UUID][]*domain.ServiceJob{carID: {j}}}
	s := NewService(&je, tCar{carID: {ID: carID, OwnerID: uuid.New()}}, tUser{empID: emp}, nil)

	_, err := s.SaveReception(context.Background(), jobID, SaveReceptionInput{OdometerKM: 100, GeneralNotes: "in"}, empID)
	require.NoError(t, err)

	_, err = s.SaveHandover(context.Background(), jobID, SaveHandoverInput{OdometerKM: 100, GeneralNotes: "out"}, empID)
	require.NoError(t, err)

	updated, err := s.jobRepo.GetByID(context.Background(), jobID)
	require.NoError(t, err)
	assert.Equal(t, domain.ServiceJobStatusClosed, updated.Status)
	assert.NotNil(t, updated.ClosedAt)
}

func TestService_SaveReception_InvalidOdometer(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	jobID := uuid.New()
	carID := uuid.New()
	emp, _ := domain.NewUser("e@t", "p", "E", "E", domain.RoleEmployee)
	emp.ID = empID
	j := &domain.ServiceJob{ID: jobID, CarID: carID, Status: domain.ServiceJobStatusOpen, OpenedByUserID: empID, OpenedAt: time.Now().UTC()}
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{jobID: j}, byCar: map[uuid.UUID][]*domain.ServiceJob{carID: {j}}}
	s := NewService(&je, tCar{carID: {ID: carID, OwnerID: uuid.New()}}, tUser{empID: emp}, nil)

	_, err := s.SaveReception(context.Background(), jobID, SaveReceptionInput{OdometerKM: -1}, empID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidServiceJobData)
}

func TestService_ListOpenedOn_EmptyDay(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	emp, _ := domain.NewUser("e@t", "p", "E", "E", domain.RoleEmployee)
	emp.ID = empID
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{}}
	s := NewService(&je, tCar{}, tUser{empID: emp}, nil)

	day := time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC)
	out, err := s.ListOpenedOn(context.Background(), day, empID)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Len(t, out, 0)
}

func TestService_ListOpenedOn_OneJob(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	carID := uuid.New()
	jobID := uuid.New()
	emp, _ := domain.NewUser("e@t", "p", "E", "E", domain.RoleEmployee)
	emp.ID = empID
	opened := time.Date(2024, 6, 10, 14, 30, 0, 0, time.UTC)
	j := &domain.ServiceJob{ID: jobID, CarID: carID, Status: domain.ServiceJobStatusOpen, OpenedByUserID: empID, OpenedAt: opened}
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{jobID: j}, byCar: map[uuid.UUID][]*domain.ServiceJob{carID: {j}}}
	s := NewService(&je, tCar{carID: {ID: carID, OwnerID: uuid.New()}}, tUser{empID: emp}, nil)

	out, err := s.ListOpenedOn(context.Background(), time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC), empID)
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, jobID, out[0].ID)
}

func TestService_ListOpenedOn_ClientDenied(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	cu, _ := domain.NewUser("c@t", "p", "C", "L", domain.RoleClient)
	cu.ID = clientID
	je := stubJobRepo{byID: map[uuid.UUID]*domain.ServiceJob{}}
	s := NewService(&je, tCar{}, tUser{clientID: cu}, nil)
	_, err := s.ListOpenedOn(context.Background(), time.Now().UTC(), clientID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
}
