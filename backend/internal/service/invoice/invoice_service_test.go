package invoice

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

type invTestUserRepo struct {
	users map[uuid.UUID]*domain.User
}

func (r *invTestUserRepo) Create(ctx context.Context, user *domain.User) error { return nil }
func (r *invTestUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (r *invTestUserRepo) GetByRole(ctx context.Context, role string, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *invTestUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *invTestUserRepo) Update(ctx context.Context, user *domain.User) error { return nil }
func (r *invTestUserRepo) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (r *invTestUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return nil
}
func (r *invTestUserRepo) GetActiveUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *invTestUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

type stubInvoiceRepo struct {
	byID      map[uuid.UUID]*domain.Invoice
	byCust    map[uuid.UUID][]*domain.Invoice
	listTotal int64
	updateErr error
}

func (s *stubInvoiceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Invoice, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}

func (s *stubInvoiceRepo) Update(ctx context.Context, invoice *domain.Invoice) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	if s.byID != nil && invoice != nil {
		cp := *invoice
		s.byID[invoice.ID] = &cp
	}
	return nil
}

func (s *stubInvoiceRepo) ListByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*domain.Invoice, int64, error) {
	if s.byCust == nil {
		return nil, 0, nil
	}
	return s.byCust[customerID], s.listTotal, nil
}

func TestInvoiceService_GetInvoice_ClientOwn(t *testing.T) {
	t.Parallel()
	cust := uuid.New()
	invID := uuid.New()
	u, err := domain.NewUser("c@x.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	u.ID = cust

	inv := &domain.Invoice{
		ID:         invID,
		CustomerID: cust,
		Amount:     100,
		Status:     "open",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	svc := NewInvoiceService(
		&stubInvoiceRepo{byID: map[uuid.UUID]*domain.Invoice{invID: inv}},
		&invTestUserRepo{users: map[uuid.UUID]*domain.User{cust: u}},
	)
	out, err := svc.GetInvoice(context.Background(), invID, cust)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, invID, out.ID)
}

func TestInvoiceService_GetInvoice_ClientOtherDenied(t *testing.T) {
	t.Parallel()
	cust := uuid.New()
	other := uuid.New()
	invID := uuid.New()
	u, err := domain.NewUser("c@x.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	u.ID = cust

	inv := &domain.Invoice{ID: invID, CustomerID: other, Amount: 50, Status: "open", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	svc := NewInvoiceService(
		&stubInvoiceRepo{byID: map[uuid.UUID]*domain.Invoice{invID: inv}},
		&invTestUserRepo{users: map[uuid.UUID]*domain.User{cust: u}},
	)
	out, err := svc.GetInvoice(context.Background(), invID, cust)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestInvoiceService_UpdateInvoice_ClientNotesOnly(t *testing.T) {
	t.Parallel()
	cust := uuid.New()
	invID := uuid.New()
	u, err := domain.NewUser("c@x.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	u.ID = cust

	existing := &domain.Invoice{
		ID: invID, CustomerID: cust, Amount: 200, Status: "open", Notes: "old",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	repo := &stubInvoiceRepo{byID: map[uuid.UUID]*domain.Invoice{invID: existing}}
	svc := NewInvoiceService(repo, &invTestUserRepo{users: map[uuid.UUID]*domain.User{cust: u}})

	out, err := svc.UpdateInvoice(context.Background(), &domain.Invoice{ID: invID, Notes: "new note", Status: "paid"}, cust)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "new note", out.Notes)
	assert.Equal(t, "open", out.Status, "client cannot change status via this path")
	assert.Equal(t, 200.0, out.Amount)
}

func TestInvoiceService_UpdateInvoice_ClientOtherDenied(t *testing.T) {
	t.Parallel()
	cust := uuid.New()
	other := uuid.New()
	invID := uuid.New()
	u, err := domain.NewUser("c@x.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	u.ID = cust

	existing := &domain.Invoice{ID: invID, CustomerID: other, Amount: 1, Status: "open", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	svc := NewInvoiceService(
		&stubInvoiceRepo{byID: map[uuid.UUID]*domain.Invoice{invID: existing}},
		&invTestUserRepo{users: map[uuid.UUID]*domain.User{cust: u}},
	)
	out, err := svc.UpdateInvoice(context.Background(), &domain.Invoice{ID: invID, Notes: "x"}, cust)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
	assert.Nil(t, out)
}

func TestInvoiceService_ListMyInvoices_ClientOnly(t *testing.T) {
	t.Parallel()
	cust := uuid.New()
	u, err := domain.NewUser("c@x.com", "pw", "C", "L", domain.RoleClient)
	require.NoError(t, err)
	u.ID = cust

	emp, err := domain.NewUser("e@x.com", "pw", "E", "E", domain.RoleEmployee)
	require.NoError(t, err)
	empID := uuid.New()
	emp.ID = empID

	repo := &stubInvoiceRepo{byCust: map[uuid.UUID][]*domain.Invoice{
		cust: {{ID: uuid.New(), CustomerID: cust, Amount: 1, Status: "open", CreatedAt: time.Now(), UpdatedAt: time.Now()}},
	}, listTotal: 1}
	svc := NewInvoiceService(repo, &invTestUserRepo{users: map[uuid.UUID]*domain.User{cust: u, empID: emp}})
	list, total, err := svc.ListMyInvoices(context.Background(), cust, 10, 0)
	require.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, int64(1), total)

	_, _, err = svc.ListMyInvoices(context.Background(), empID, 10, 0)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
}
