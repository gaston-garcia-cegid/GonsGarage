package repair

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type repairTestUserRepo struct {
	users  map[uuid.UUID]*domain.User
	getErr error
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
func (r *repairTestUserRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (r *repairTestUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return nil
}
func (r *repairTestUserRepo) GetActiveUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}

func (r *repairTestUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	if r.users == nil {
		return nil, nil
	}
	u, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

type stubRepairCarRepo struct {
	byID map[uuid.UUID]*domain.Car
}

func (s *stubRepairCarRepo) Create(ctx context.Context, car *domain.Car) error { return nil }
func (s *stubRepairCarRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	if s.byID == nil {
		return nil, domain.ErrCarNotFound
	}
	c, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrCarNotFound
	}
	return c, nil
}
func (s *stubRepairCarRepo) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Car, error) {
	return nil, nil
}
func (s *stubRepairCarRepo) GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	return nil, nil
}
func (s *stubRepairCarRepo) List(ctx context.Context, limit, offset int) ([]*domain.Car, error) { return nil, nil }
func (s *stubRepairCarRepo) Update(ctx context.Context, car *domain.Car) error { return nil }
func (s *stubRepairCarRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (s *stubRepairCarRepo) GetWithRepairs(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	return nil, nil
}
func (s *stubRepairCarRepo) GetDeletedByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error) {
	return nil, nil
}
func (s *stubRepairCarRepo) Restore(ctx context.Context, id uuid.UUID) error { return nil }

type stubRepairRepo struct {
	byID  map[uuid.UUID]*domain.Repair
	byCar map[uuid.UUID][]*domain.Repair
}

func newStubRepairRepo() *stubRepairRepo {
	return &stubRepairRepo{
		byID:  make(map[uuid.UUID]*domain.Repair),
		byCar: make(map[uuid.UUID][]*domain.Repair),
	}
}

func (s *stubRepairRepo) Create(ctx context.Context, r *domain.Repair) error {
	cp := *r
	s.byID[r.ID] = &cp
	s.byCar[r.CarID] = append(s.byCar[r.CarID], &cp)
	return nil
}

func (s *stubRepairRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error) {
	r, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrRepairNotFound
	}
	return r, nil
}

func (s *stubRepairRepo) GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error) {
	return s.byCar[carID], nil
}

func (s *stubRepairRepo) Update(ctx context.Context, r *domain.Repair) error {
	if _, ok := s.byID[r.ID]; !ok {
		return domain.ErrRepairNotFound
	}
	cp := *r
	s.byID[r.ID] = &cp
	return nil
}

func (s *stubRepairRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }

func TestRepairService_CreateRepair_ClientForbidden(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	carID := uuid.New()
	svc := NewRepairService(
		newStubRepairRepo(),
		&stubRepairCarRepo{byID: map[uuid.UUID]*domain.Car{carID: {ID: carID, OwnerID: clientID}}},
		&repairTestUserRepo{users: map[uuid.UUID]*domain.User{
			clientID: {ID: clientID, Role: domain.RoleClient},
		}},
	)
	out, err := svc.CreateRepair(context.Background(), &domain.Repair{CarID: carID, Description: "Oil"}, clientID)
	require.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestRepairService_CreateRepair_EmployeeSetsTechnician(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	clientID := uuid.New()
	carID := uuid.New()
	repRepo := newStubRepairRepo()
	svc := NewRepairService(
		repRepo,
		&stubRepairCarRepo{byID: map[uuid.UUID]*domain.Car{carID: {ID: carID, OwnerID: clientID}}},
		&repairTestUserRepo{users: map[uuid.UUID]*domain.User{
			empID: {ID: empID, Role: domain.RoleEmployee},
		}},
	)
	out, err := svc.CreateRepair(context.Background(), &domain.Repair{CarID: carID, Description: "Brakes", Cost: 120}, empID)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, empID, out.TechnicianID)
	assert.Equal(t, domain.RepairStatusPending, out.Status)
	assert.Len(t, repRepo.byCar[carID], 1)
}

func TestRepairService_CreateRepair_CarNotFound(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	svc := NewRepairService(
		newStubRepairRepo(),
		&stubRepairCarRepo{byID: map[uuid.UUID]*domain.Car{}},
		&repairTestUserRepo{users: map[uuid.UUID]*domain.User{
			empID: {ID: empID, Role: domain.RoleEmployee},
		}},
	)
	out, err := svc.CreateRepair(context.Background(), &domain.Repair{CarID: uuid.New(), Description: "x"}, empID)
	require.ErrorIs(t, err, domain.ErrCarNotFound)
	assert.Nil(t, out)
}

func TestRepairService_GetRepairsByCarID_ClientWrongOwner(t *testing.T) {
	t.Parallel()
	clientA := uuid.New()
	clientB := uuid.New()
	carID := uuid.New()
	svc := NewRepairService(
		newStubRepairRepo(),
		&stubRepairCarRepo{byID: map[uuid.UUID]*domain.Car{carID: {ID: carID, OwnerID: clientA}}},
		&repairTestUserRepo{users: map[uuid.UUID]*domain.User{
			clientB: {ID: clientB, Role: domain.RoleClient},
		}},
	)
	out, err := svc.GetRepairsByCarID(context.Background(), carID, clientB)
	require.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestRepairService_GetRepairsByCarID_ClientSeesOwn(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	empID := uuid.New()
	carID := uuid.New()
	repRepo := newStubRepairRepo()
	carRepo := &stubRepairCarRepo{byID: map[uuid.UUID]*domain.Car{carID: {ID: carID, OwnerID: clientID}}}
	users := &repairTestUserRepo{users: map[uuid.UUID]*domain.User{
		empID:     {ID: empID, Role: domain.RoleEmployee},
		clientID:  {ID: clientID, Role: domain.RoleClient},
	}}
	svcEmp := NewRepairService(repRepo, carRepo, users)
	_, err := svcEmp.CreateRepair(context.Background(), &domain.Repair{CarID: carID, Description: "Tires"}, empID)
	require.NoError(t, err)

	svcClient := NewRepairService(repRepo, carRepo, users)
	list, err := svcClient.GetRepairsByCarID(context.Background(), carID, clientID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, "Tires", list[0].Description)
}

func TestRepairService_UpdateRepair_ClientForbidden(t *testing.T) {
	t.Parallel()
	clientID := uuid.New()
	rid := uuid.New()
	svc := NewRepairService(
		newStubRepairRepo(),
		&stubRepairCarRepo{},
		&repairTestUserRepo{users: map[uuid.UUID]*domain.User{
			clientID: {ID: clientID, Role: domain.RoleClient},
		}},
	)
	out, err := svc.UpdateRepair(context.Background(), &domain.Repair{ID: rid, Status: domain.RepairStatusCompleted}, clientID)
	require.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestRepairService_UpdateRepair_EmployeeCompletes(t *testing.T) {
	t.Parallel()
	empID := uuid.New()
	carID := uuid.New()
	rid := uuid.New()
	now := time.Now().UTC()
	repRepo := newStubRepairRepo()
	require.NoError(t, repRepo.Create(context.Background(), &domain.Repair{
		ID: rid, CarID: carID, TechnicianID: empID, Description: "Align",
		Status: domain.RepairStatusInProgress, Cost: 50, CreatedAt: now, UpdatedAt: now,
	}))
	svc := NewRepairService(
		repRepo,
		&stubRepairCarRepo{byID: map[uuid.UUID]*domain.Car{carID: {ID: carID}}},
		&repairTestUserRepo{users: map[uuid.UUID]*domain.User{
			empID: {ID: empID, Role: domain.RoleEmployee},
		}},
	)
	out, err := svc.UpdateRepair(context.Background(), &domain.Repair{
		ID: rid, Status: domain.RepairStatusCompleted, Cost: 50,
	}, empID)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, domain.RepairStatusCompleted, out.Status)
	assert.NotNil(t, out.CompletedAt)
}

func TestRepairService_GetRepair_UserRepoError(t *testing.T) {
	t.Parallel()
	svc := NewRepairService(newStubRepairRepo(), &stubRepairCarRepo{}, &repairTestUserRepo{getErr: errors.New("db")})
	out, err := svc.GetRepair(context.Background(), uuid.New(), uuid.New())
	require.Error(t, err)
	assert.Nil(t, out)
}
