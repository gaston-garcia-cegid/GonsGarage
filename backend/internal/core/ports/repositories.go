package ports

import (
	"context"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

// UserRepository define os métodos para o repositório de usuários
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*domain.User, error)
}

// EmployeeRepository define os métodos para o repositório de funcionários
type EmployeeRepository interface {
	Create(ctx context.Context, employee *domain.Employee) error
	FindByID(ctx context.Context, id uint) (*domain.Employee, error)
	Update(ctx context.Context, employee *domain.Employee) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filters *EmployeeFilters) ([]*domain.Employee, int64, error)
}

// CacheRepository define os métodos para o repositório de cache
type CacheRepository interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Delete(ctx context.Context, key string) error
}

// SessionRepository define os métodos para o repositório de sessões
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id string) (*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*domain.Session, int64, error)
	SearchByName(ctx context.Context, name string, limit int) ([]*domain.Session, error)
}

// QueueRepository define os métodos para o repositório de filas
type QueueRepository interface {
	Enqueue(ctx context.Context, job *domain.Job) error
	Dequeue(ctx context.Context) (*domain.Job, error)
	Acknowledge(ctx context.Context, id string) error
	Reject(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*domain.Job, int64, error)
	SearchByName(ctx context.Context, name string, limit int) ([]*domain.Job, error)
}

// WorkshopRepository defines the interface for the workshop repository
type WorkshopRepository interface {
	CreateWorkshop(ctx context.Context, workshop *domain.Workshop) error
	UpdateWorkshop(ctx context.Context, workshop *domain.Workshop) error
	DeleteWorkshop(ctx context.Context, id string) error
	GetWorkshopByID(ctx context.Context, id string) (*domain.Workshop, error)
	ListWorkshops(ctx context.Context, limit, offset int) ([]*domain.Workshop, int64, error)
}

// NotificationRepository defines the interface for the notification repository
type NotificationRepository interface {
	SendWorkshopCreatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopUpdatedNotification(ctx context.Context, workshop *domain.Workshop) error
	SendWorkshopDeletedNotification(ctx context.Context, workshop *domain.Workshop) error
}

// AccountingRepository defines the interface for the accounting repository
type AccountingRepository interface {
	CreateAccountingEntry(ctx context.Context, entry *domain.AccountingEntry) error
	GetAccountingEntryByID(ctx context.Context, id string) (*domain.AccountingEntry, error)
	UpdateAccountingEntry(ctx context.Context, entry *domain.AccountingEntry) error
	DeleteAccountingEntry(ctx context.Context, id string) error
	ListAccountingEntries(ctx context.Context, limit, offset int) ([]*domain.AccountingEntry, int64, error)
}
