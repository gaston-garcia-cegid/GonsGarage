package repair

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

type repairTestUserRepo struct {
	users map[uuid.UUID]*domain.User
}

func (r *repairTestUserRepo) Create(ctx context.Context, user *domain.User) error { return nil }
func (r *repairTestUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (r *repairTestUserRepo) GetByRole(ctx context.Context, role string, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *repairTestUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *repairTestUserRepo) Update(ctx context.Context, user *domain.User) error { return nil }
func (r *repairTestUserRepo) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (r *repairTestUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return nil
}
func (r *repairTestUserRepo) GetActiveUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *repairTestUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

type stubRepairRepo struct {
	created []*domain.Repair
	byID    map[uuid.UUID]*domain.Repair
	byCar   map[uuid.UUID][]*domain.Repair
	deleted []uuid.UUID
}

func (s *stubRepairRepo) Create(ctx context.Context, repair *domain.Repair) error {
	s.created = append(s.created, repair)
	return nil
}

func (s *stubRepairRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error) {
	if s.byID == nil {
		return nil, errors.New("not found")
	}
	r, ok := s.byID[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return r, nil
}

func (s *stubRepairRepo) Update(ctx context.Context, repair *domain.Repair) error { return nil }
func (s *stubRepairRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if s.byID != nil {
		if _, ok := s.byID[id]; !ok {
			return domain.ErrRepairNotFound
		}
		delete(s.byID, id)
	}
	s.deleted = append(s.deleted, id)
	return nil
}

func (s *stubRepairRepo) GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error) {
	if s.byCar == nil {
		return nil, nil
	}
	return s.byCar[carID], nil
}

func (s *stubRepairRepo) ListIDsByServiceJobID(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

type repairStubCarRepo struct {
	byID map[uuid.UUID]*domain.Car
}

func (s *repairStubCarRepo) Create(ctx context.Context, car *domain.Car) error { return nil }
func (s *repairStubCarRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	if s.byID == nil {
		return nil, domain.ErrCarNotFound
	}
	c, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrCarNotFound
	}
	return c, nil
}
func (s *repairStubCarRepo) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Car, error) {
	return nil, nil
}
func (s *repairStubCarRepo) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	return nil, nil
}
func (s *repairStubCarRepo) List(ctx context.Context, limit, offset int) ([]*domain.Car, error) {
	return nil, nil
}
func (s *repairStubCarRepo) Update(ctx context.Context, car *domain.Car) error { return nil }
func (s *repairStubCarRepo) Delete(ctx context.Context, id uuid.UUID) error    { return nil }
func (s *repairStubCarRepo) GetWithRepairs(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	return nil, nil
}
func (s *repairStubCarRepo) GetDeletedByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	return nil, nil
}
func (s *repairStubCarRepo) Restore(ctx context.Context, id uuid.UUID) error { return nil }

func TestRepairService_CreateRepair_EmployeeOnClientCar(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	empID := uuid.New()
	carID := uuid.New()

	client, err := domain.NewUser("client@example.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	client.ID = clientID

	emp, err := domain.NewUser("mech@example.com", "pw", "M", "E", domain.RoleEmployee)
	require.NoError(t, err)
	emp.ID = empID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{
		empID:    emp,
		clientID: client,
	}}
	carRepo := &repairStubCarRepo{byID: map[uuid.UUID]*domain.Car{
		carID: {ID: carID, OwnerID: clientID, LicensePlate: "R-1"},
	}}
	repairRepo := &stubRepairRepo{byCar: map[uuid.UUID][]*domain.Repair{carID: {}}}

	svc := NewRepairService(repairRepo, carRepo, userRepo)
	now := time.Now().UTC()
	repair := &domain.Repair{
		CarID:       carID,
		Description: "Oil change",
		Status:      domain.RepairStatusPending,
		Cost:        50,
		StartedAt:   &now,
	}

	out, err := svc.CreateRepair(context.Background(), repair, empID)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, empID, out.TechnicianID)
	assert.Equal(t, carID, out.CarID)
	assert.NotEqual(t, uuid.Nil, out.ID)
	require.Len(t, repairRepo.created, 1)
}

func TestRepairService_CreateRepair_ClientDenied(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	carID := uuid.New()
	client, err := domain.NewUser("c@example.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	client.ID = clientID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{clientID: client}}
	carRepo := &repairStubCarRepo{byID: map[uuid.UUID]*domain.Car{
		carID: {ID: carID, OwnerID: clientID},
	}}
	svc := NewRepairService(&stubRepairRepo{}, carRepo, userRepo)

	out, err := svc.CreateRepair(context.Background(), &domain.Repair{
		CarID: carID, Description: "x", Status: domain.RepairStatusPending, Cost: 1,
	}, clientID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestRepairService_GetRepair_ClientOwnCar(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	carID := uuid.New()
	repairID := uuid.New()

	client, err := domain.NewUser("c@example.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	client.ID = clientID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{clientID: client}}
	carRepo := &repairStubCarRepo{byID: map[uuid.UUID]*domain.Car{
		carID: {ID: carID, OwnerID: clientID},
	}}
	repairRepo := &stubRepairRepo{byID: map[uuid.UUID]*domain.Repair{
		repairID: {ID: repairID, CarID: carID, Description: "Brakes", Status: domain.RepairStatusPending, Cost: 100},
	}}

	svc := NewRepairService(repairRepo, carRepo, userRepo)
	out, err := svc.GetRepair(context.Background(), repairID, clientID)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, repairID, out.ID)
}

func TestRepairService_GetRepair_ClientOtherOwnersCarDenied(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	otherOwner := uuid.New()
	carID := uuid.New()
	repairID := uuid.New()

	client, err := domain.NewUser("c@example.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	client.ID = clientID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{clientID: client}}
	carRepo := &repairStubCarRepo{byID: map[uuid.UUID]*domain.Car{
		carID: {ID: carID, OwnerID: otherOwner},
	}}
	repairRepo := &stubRepairRepo{byID: map[uuid.UUID]*domain.Repair{
		repairID: {ID: repairID, CarID: carID, Description: "x", Status: domain.RepairStatusPending},
	}}

	svc := NewRepairService(repairRepo, carRepo, userRepo)
	out, err := svc.GetRepair(context.Background(), repairID, clientID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestRepairService_GetRepairsByCarID_ClientOtherCarDenied(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	other := uuid.New()
	carID := uuid.New()
	client, err := domain.NewUser("c@example.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	client.ID = clientID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{clientID: client}}
	carRepo := &repairStubCarRepo{byID: map[uuid.UUID]*domain.Car{
		carID: {ID: carID, OwnerID: other},
	}}
	svc := NewRepairService(&stubRepairRepo{}, carRepo, userRepo)

	out, err := svc.GetRepairsByCarID(context.Background(), carID, clientID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestRepairService_UpdateRepair_ClientDenied(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	client, err := domain.NewUser("c@example.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	client.ID = clientID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{clientID: client}}
	svc := NewRepairService(&stubRepairRepo{}, &repairStubCarRepo{}, userRepo)

	out, err := svc.UpdateRepair(context.Background(), &domain.Repair{
		ID: uuid.New(), CarID: uuid.New(), Description: "x", Status: domain.RepairStatusPending, Cost: 1,
	}, clientID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestRepairService_DeleteRepair_ClientDenied(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	client, err := domain.NewUser("c@example.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	client.ID = clientID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{clientID: client}}
	svc := NewRepairService(&stubRepairRepo{}, &repairStubCarRepo{}, userRepo)

	err = svc.DeleteRepair(context.Background(), uuid.New(), clientID)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
}

func TestRepairService_DeleteRepair_EmployeeOK(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	repairID := uuid.New()
	carID := uuid.New()

	emp, err := domain.NewUser("mech@example.com", "pw", "M", "E", domain.RoleEmployee)
	require.NoError(t, err)
	emp.ID = empID

	userRepo := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{empID: emp}}
	repairRepo := &stubRepairRepo{byID: map[uuid.UUID]*domain.Repair{
		repairID: {ID: repairID, CarID: carID, Description: "x", Status: domain.RepairStatusPending, Cost: 1},
	}}
	svc := NewRepairService(repairRepo, &repairStubCarRepo{}, userRepo)

	err = svc.DeleteRepair(context.Background(), repairID, empID)
	require.NoError(t, err)
	require.Contains(t, repairRepo.deleted, repairID)
}
