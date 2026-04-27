package part

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/ports"
	"github.com/gaston-garcia-cegid/gonsgarage/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- test doubles ---

type partTestUserRepo struct {
	users map[uuid.UUID]*domain.User
}

func (r *partTestUserRepo) Create(ctx context.Context, user *domain.User) error { return nil }
func (r *partTestUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}
func (r *partTestUserRepo) GetByRole(ctx context.Context, role string, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *partTestUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) { return nil, nil }
func (r *partTestUserRepo) Update(ctx context.Context, user *domain.User) error { return nil }
func (r *partTestUserRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (r *partTestUserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return nil
}
func (r *partTestUserRepo) GetActiveUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}
func (r *partTestUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

type stubPartRepo struct {
	byID      map[uuid.UUID]*domain.PartItem
	byBarcode map[string]*domain.PartItem
}

func newStubPartRepo() *stubPartRepo {
	return &stubPartRepo{
		byID:      make(map[uuid.UUID]*domain.PartItem),
		byBarcode: make(map[string]*domain.PartItem),
	}
}

func (s *stubPartRepo) indexBarcode(p *domain.PartItem) {
	if p == nil {
		return
	}
	b := p.Barcode
	if b == "" {
		return
	}
	cp := *p
	s.byBarcode[b] = &cp
}

func (s *stubPartRepo) Create(ctx context.Context, p *domain.PartItem) error {
	if p == nil {
		return errors.New("nil")
	}
	cp := *p
	s.byID[p.ID] = &cp
	s.indexBarcode(&cp)
	return nil
}

func (s *stubPartRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.PartItem, error) {
	v, ok := s.byID[id]
	if !ok {
		return nil, domain.ErrPartItemNotFound
	}
	cp := *v
	return &cp, nil
}

func (s *stubPartRepo) GetByBarcode(ctx context.Context, barcode string) (*domain.PartItem, error) {
	if barcode == "" {
		return nil, domain.ErrPartItemNotFound
	}
	v, ok := s.byBarcode[barcode]
	if !ok {
		return nil, domain.ErrPartItemNotFound
	}
	cp := *v
	return &cp, nil
}

func (s *stubPartRepo) Update(ctx context.Context, p *domain.PartItem) error {
	if _, ok := s.byID[p.ID]; !ok {
		return domain.ErrPartItemNotFound
	}
	old := s.byID[p.ID]
	if old != nil && old.Barcode != "" {
		delete(s.byBarcode, old.Barcode)
	}
	cp := *p
	s.byID[p.ID] = &cp
	s.indexBarcode(&cp)
	return nil
}

func (s *stubPartRepo) Delete(ctx context.Context, id uuid.UUID) error {
	v, ok := s.byID[id]
	if !ok {
		return domain.ErrPartItemNotFound
	}
	if v.Barcode != "" {
		delete(s.byBarcode, v.Barcode)
	}
	delete(s.byID, id)
	return nil
}

func (s *stubPartRepo) List(ctx context.Context, f ports.PartItemListFilters) ([]*domain.PartItem, int64, error) {
	var out []*domain.PartItem
	for _, v := range s.byID {
		if v == nil {
			continue
		}
		cp := *v
		out = append(out, &cp)
	}
	return out, int64(len(out)), nil
}

func managerUser(id uuid.UUID) *domain.User {
	u, _ := domain.NewUser("m@x.com", "pw", "M", "M", domain.RoleManager)
	u.ID = id
	return u
}

func employeeUser(id uuid.UUID) *domain.User {
	u, _ := domain.NewUser("e@x.com", "pw", "E", "E", domain.RoleEmployee)
	u.ID = id
	return u
}

func adminUser(id uuid.UUID) *domain.User {
	u, _ := domain.NewUser("a@x.com", "pw", "A", "A", domain.RoleAdmin)
	u.ID = id
	return u
}

func basePart() *domain.PartItem {
	id := uuid.New()
	return &domain.PartItem{
		ID:        id,
		Reference: "R1",
		Brand:     "B",
		Name:      "N",
		Barcode:   "BC-1",
		Quantity:  1,
		UOM:       domain.PartUOMUnit,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestPartService_Create_negativeQuantity(t *testing.T) {
	t.Parallel()
	mid := uuid.New()
	svc := NewPartService(newStubPartRepo(), &partTestUserRepo{users: map[uuid.UUID]*domain.User{mid: managerUser(mid)}})
	p := basePart()
	p.Quantity = -1
	_, err := svc.Create(context.Background(), p, mid)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "quantity")
}

func TestPartService_Create_invalidUOM(t *testing.T) {
	t.Parallel()
	mid := uuid.New()
	svc := NewPartService(newStubPartRepo(), &partTestUserRepo{users: map[uuid.UUID]*domain.User{mid: managerUser(mid)}})
	p := basePart()
	p.UOM = "kg"
	_, err := svc.Create(context.Background(), p, mid)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "uom")
}

func TestPartService_Create_duplicateBarcode(t *testing.T) {
	t.Parallel()
	mid := uuid.New()
	repo := newStubPartRepo()
	existing := basePart()
	existing.Barcode = "SHARED"
	require.NoError(t, repo.Create(context.Background(), existing))

	svc := NewPartService(repo, &partTestUserRepo{users: map[uuid.UUID]*domain.User{mid: managerUser(mid)}})
	dup := basePart()
	dup.ID = uuid.New()
	dup.Barcode = "SHARED"
	_, err := svc.Create(context.Background(), dup, mid)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrPartItemDuplicateBarcode)
}

func TestPartService_Update_duplicateBarcode(t *testing.T) {
	t.Parallel()
	mid := uuid.New()
	repo := newStubPartRepo()
	a := basePart()
	a.Barcode = "A1"
	b := basePart()
	b.ID = uuid.New()
	b.Barcode = "B1"
	require.NoError(t, repo.Create(context.Background(), a))
	require.NoError(t, repo.Create(context.Background(), b))

	svc := NewPartService(repo, &partTestUserRepo{users: map[uuid.UUID]*domain.User{mid: managerUser(mid)}})
	b.Barcode = "A1"
	_, err := svc.Update(context.Background(), b, mid)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrPartItemDuplicateBarcode)
}

func TestPartService_Create_negativeMinimumQuantity(t *testing.T) {
	t.Parallel()
	mid := uuid.New()
	svc := NewPartService(newStubPartRepo(), &partTestUserRepo{users: map[uuid.UUID]*domain.User{mid: managerUser(mid)}})
	p := basePart()
	bad := -3.0
	p.MinimumQuantity = &bad
	_, err := svc.Create(context.Background(), p, mid)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "minimum")
}

func TestPartService_Create_managerOK(t *testing.T) {
	t.Parallel()
	mid := uuid.New()
	svc := NewPartService(newStubPartRepo(), &partTestUserRepo{users: map[uuid.UUID]*domain.User{mid: managerUser(mid)}})
	p := basePart()
	out, err := svc.Create(context.Background(), p, mid)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, p.ID, out.ID)
	assert.Equal(t, "BC-1", out.Barcode)
}

func TestPartService_Create_adminOK(t *testing.T) {
	t.Parallel()
	aid := uuid.New()
	svc := NewPartService(newStubPartRepo(), &partTestUserRepo{users: map[uuid.UUID]*domain.User{aid: adminUser(aid)}})
	p := basePart()
	p.ID = uuid.New()
	p.Barcode = "ADM-BC"
	out, err := svc.Create(context.Background(), p, aid)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "ADM-BC", out.Barcode)
}

func TestPartService_Update_sameBarcodeAllowed(t *testing.T) {
	t.Parallel()
	mid := uuid.New()
	repo := newStubPartRepo()
	p := basePart()
	require.NoError(t, repo.Create(context.Background(), p))
	svc := NewPartService(repo, &partTestUserRepo{users: map[uuid.UUID]*domain.User{mid: managerUser(mid)}})
	p.Name = "Renamed"
	out, err := svc.Update(context.Background(), p, mid)
	require.NoError(t, err)
	require.NotNil(t, out)
	assert.Equal(t, "Renamed", out.Name)
	assert.Equal(t, "BC-1", out.Barcode)
}

func TestPartService_Create_employeeForbidden(t *testing.T) {
	t.Parallel()
	eid := uuid.New()
	svc := NewPartService(newStubPartRepo(), &partTestUserRepo{users: map[uuid.UUID]*domain.User{eid: employeeUser(eid)}})
	p := basePart()
	_, err := svc.Create(context.Background(), p, eid)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrUnauthorizedAccess)
}
