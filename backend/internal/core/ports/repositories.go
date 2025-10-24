package ports

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/gonsgarage/internal/core/domain"
	"github.com/google/uuid"
)

// UserRepository define os métodos para o repositório de usuários
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByRole(ctx context.Context, role string, limit, offset int) ([]*domain.User, error)
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	GetActiveUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

// EmployeeRepository define os métodos para o repositório de funcionários
type EmployeeRepository interface {
	Create(ctx context.Context, employee *domain.Employee) error
	FindByID(ctx context.Context, id uint) (*domain.Employee, error)
	Update(ctx context.Context, employee *domain.Employee) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filters *EmployeeFilters) ([]*domain.Employee, int64, error)
}

// EmployeeFilters representa os filtros para listagem de funcionários
type EmployeeFilters struct {
	Department *string
	IsActive   *bool
	Search     *string
	SortBy     string
	SortOrder  string
	Limit      int
	Offset     int
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

// CarRepository defines the interface for the car repository
type CarRepository interface {
	Create(ctx context.Context, car *domain.Car) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Car, error)
	GetByLicensePlate(ctx context.Context, licensePlate string) (*domain.Car, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Car, error)
	Update(ctx context.Context, car *domain.Car) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetWithRepairs(ctx context.Context, id uuid.UUID) (*domain.Car, error)
}

// RepairRepository defines the interface for the repair repository
type RepairRepository interface {
	Create(ctx context.Context, repair *domain.Repair) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Repair, error)
	Update(ctx context.Context, repair *domain.Repair) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByCarID(ctx context.Context, carID uuid.UUID) ([]*domain.Repair, error)
}

// Logger defines the interface for logging
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type ClientRepository interface {
	Create(ctx context.Context, client *domain.Client) error
	GetByEmail(ctx context.Context, email string) (*domain.Client, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error)
	GetWithCars(ctx context.Context, id uuid.UUID) (*domain.Client, error)
	List(ctx context.Context) ([]*domain.Client, error)
	Search(ctx context.Context, name string, limit int) ([]*domain.Client, error)
	Update(ctx context.Context, client *domain.Client) (*domain.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context) (int64, error)
	GetActiveClients(ctx context.Context) ([]*domain.Client, error)
	DeactivateClient(ctx context.Context, id uuid.UUID) error
	ActivateClient(ctx context.Context, id uuid.UUID) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetClientsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*domain.Client, error)
}

// ClientUseCase defines the business logic interface for clients
type ClientUseCase interface {
	CreateClient(ctx context.Context, client *domain.Client) (*domain.Client, error)
	GetClient(ctx context.Context, clientID uuid.UUID, requestingUserID uuid.UUID) (*domain.Client, error)
	GetClientProfile(ctx context.Context, clientID uuid.UUID) (*domain.Client, error)
	UpdateClient(ctx context.Context, client *domain.Client, requestingUserID uuid.UUID) (*domain.Client, error)
	UpdateClientProfile(ctx context.Context, clientID uuid.UUID, client *domain.Client) (*domain.Client, error)
	DeleteClient(ctx context.Context, clientID uuid.UUID) error
	ListClients(ctx context.Context) ([]*domain.Client, error)
	GetClientCars(ctx context.Context, clientID uuid.UUID, requestingUserID uuid.UUID) ([]*domain.Car, error)
	GetClientRepairs(ctx context.Context, clientID uuid.UUID, requestingUserID uuid.UUID) ([]*domain.Repair, error)
}
