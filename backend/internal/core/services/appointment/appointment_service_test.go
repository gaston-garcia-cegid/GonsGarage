package appointment

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type apptTestUserRepo struct {
	users  map[uuid.UUID]*domain.User
	getErr error
}

func (r *apptTestUserRepo) Create(ctx context.Context, user *domain.User) error { return nil }
func (r *apptTestUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (r *apptTestUserRepo) GetByRole(ctx context.Context, role string, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *apptTestUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *apptTestUserRepo) Update(ctx context.Context, user *domain.User) error { return nil }
func (r *apptTestUserRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (r *apptTestUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return nil
}
func (r *apptTestUserRepo) GetActiveUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}

func (r *apptTestUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
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

type stubApptRepo struct {
	created     []*domain.Appointment
	createErr   error
	byID        map[uuid.UUID]*domain.Appointment
	lastList    *ports.AppointmentFilters
	listApps    []*domain.Appointment
	listTotal   int64
	listErr     error
	updateErr   error
	deleteErr   error
}

func (s *stubApptRepo) Create(ctx context.Context, a *domain.Appointment) error {
	if s.createErr != nil {
		return s.createErr
	}
	s.created = append(s.created, a)
	return nil
}

func (s *stubApptRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	if s.byID == nil {
		return nil, nil
	}
	a, ok := s.byID[id]
	if !ok {
		return nil, nil
	}
	return a, nil
}

func (s *stubApptRepo) Update(ctx context.Context, a *domain.Appointment) error {
	return s.updateErr
}

func (s *stubApptRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return s.deleteErr
}

func (s *stubApptRepo) List(ctx context.Context, filters *ports.AppointmentFilters) ([]*domain.Appointment, int64, error) {
	s.lastList = filters
	if s.listErr != nil {
		return nil, 0, s.listErr
	}
	return s.listApps, s.listTotal, nil
}

type noopCache struct{}

func (noopCache) Get(ctx context.Context, key string, dest interface{}) error { return nil }

func (noopCache) Set(ctx context.Context, key string, value interface{}, ttl int) error { return nil }

func (noopCache) Delete(ctx context.Context, key string) error { return nil }

func TestAppointmentService_CreateAppointment_UserNotFound(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	svc := NewAppointmentService(&stubApptRepo{}, &apptTestUserRepo{users: map[uuid.UUID]*domain.User{}}, noopCache{})
	appt := sampleAppointment(userID)

	out, err := svc.CreateAppointment(context.Background(), appt, userID)
	require.ErrorIs(t, err, domain.ErrUserNotFound)
	assert.Nil(t, out)
}

func TestAppointmentService_CreateAppointment_GetUserError(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	repoErr := errors.New("db down")
	svc := NewAppointmentService(
		&stubApptRepo{},
		&apptTestUserRepo{getErr: repoErr},
		noopCache{},
	)
	appt := sampleAppointment(userID)

	out, err := svc.CreateAppointment(context.Background(), appt, userID)
	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, out)
}

func TestAppointmentService_CreateAppointment_Success(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	user, err := domain.NewUser("u@example.com", "pw", "U", "Ser", domain.RoleClient)
	require.NoError(t, err)
	user.ID = userID

	apptRepo := &stubApptRepo{}
	svc := NewAppointmentService(
		apptRepo,
		&apptTestUserRepo{users: map[uuid.UUID]*domain.User{userID: user}},
		noopCache{},
	)
	in := sampleAppointment(userID)
	in.ID = uuid.Nil // service assigns new ID

	out, err := svc.CreateAppointment(context.Background(), in, userID)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.NotEqual(t, uuid.Nil, out.ID)
	assert.False(t, out.CreatedAt.IsZero())
	assert.False(t, out.UpdatedAt.IsZero())
	require.Len(t, apptRepo.created, 1)
	assert.Equal(t, out.ID, apptRepo.created[0].ID)
}

func TestAppointmentService_CreateAppointment_RepoError(t *testing.T) {
	t.Parallel()
	userID := uuid.New()
	user, err := domain.NewUser("u@example.com", "pw", "U", "Ser", domain.RoleClient)
	require.NoError(t, err)
	user.ID = userID

	createErr := errors.New("insert failed")
	svc := NewAppointmentService(
		&stubApptRepo{createErr: createErr},
		&apptTestUserRepo{users: map[uuid.UUID]*domain.User{userID: user}},
		noopCache{},
	)

	out, err := svc.CreateAppointment(context.Background(), sampleAppointment(userID), userID)
	require.ErrorIs(t, err, createErr)
	assert.Nil(t, out)
}

func TestAppointmentService_ListAppointments_DefaultFilters(t *testing.T) {
	t.Parallel()
	apptRepo := &stubApptRepo{listApps: []*domain.Appointment{}, listTotal: 0}
	svc := NewAppointmentService(apptRepo, &apptTestUserRepo{}, noopCache{})

	apps, total, err := svc.ListAppointments(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, apps)
	assert.Equal(t, int64(0), total)
	require.NotNil(t, apptRepo.lastList)
	assert.Equal(t, 10, apptRepo.lastList.Limit)
	assert.Equal(t, 0, apptRepo.lastList.Offset)
	assert.Equal(t, "created_at", apptRepo.lastList.SortBy)
	assert.Equal(t, "DESC", apptRepo.lastList.SortOrder)
}

func TestAppointmentService_ListAppointments_FillsEmptySortAndLimit(t *testing.T) {
	t.Parallel()
	apptRepo := &stubApptRepo{}
	svc := NewAppointmentService(apptRepo, &apptTestUserRepo{}, noopCache{})

	f := &ports.AppointmentFilters{Limit: 0, SortBy: "", SortOrder: ""}
	_, _, err := svc.ListAppointments(context.Background(), f)
	require.NoError(t, err)
	require.NotNil(t, apptRepo.lastList)
	assert.Equal(t, 10, apptRepo.lastList.Limit)
	assert.Equal(t, "created_at", apptRepo.lastList.SortBy)
	assert.Equal(t, "DESC", apptRepo.lastList.SortOrder)
}

func TestAppointmentService_GetAppointment(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	expected := &domain.Appointment{ID: id, ServiceType: "oil"}
	svc := NewAppointmentService(
		&stubApptRepo{byID: map[uuid.UUID]*domain.Appointment{id: expected}},
		&apptTestUserRepo{},
		noopCache{},
	)

	out, err := svc.GetAppointment(context.Background(), id, uuid.New())
	require.NoError(t, err)
	assert.Equal(t, expected, out)
}

func TestAppointmentService_UpdateAppointment_RepoError(t *testing.T) {
	t.Parallel()
	updErr := errors.New("update failed")
	svc := NewAppointmentService(&stubApptRepo{updateErr: updErr}, &apptTestUserRepo{}, noopCache{})
	appt := &domain.Appointment{ID: uuid.New()}

	out, err := svc.UpdateAppointment(context.Background(), appt, uuid.New())
	require.ErrorIs(t, err, updErr)
	assert.Nil(t, out)
}

func TestAppointmentService_DeleteAppointment(t *testing.T) {
	t.Parallel()
	svc := NewAppointmentService(&stubApptRepo{}, &apptTestUserRepo{}, noopCache{})
	err := svc.DeleteAppointment(context.Background(), uuid.New(), uuid.New())
	require.NoError(t, err)
}

func TestAppointmentService_DeleteAppointment_RepoError(t *testing.T) {
	t.Parallel()
	delErr := errors.New("delete failed")
	svc := NewAppointmentService(&stubApptRepo{deleteErr: delErr}, &apptTestUserRepo{}, noopCache{})

	err := svc.DeleteAppointment(context.Background(), uuid.New(), uuid.New())
	require.ErrorIs(t, err, delErr)
}

func sampleAppointment(customerID uuid.UUID) *domain.Appointment {
	carID := uuid.New()
	return &domain.Appointment{
		CustomerID:  customerID,
		CarID:       carID,
		ServiceType: "inspection",
		Status:      domain.AppointmentStatusScheduled,
		ScheduledAt: time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second),
	}
}
