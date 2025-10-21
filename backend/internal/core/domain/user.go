package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound is returned when a user is not found in the repository.
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidRole       = errors.New("invalid role")
)

type User struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Email        string     `json:"email" gorm:"unique;not null"`
	PasswordHash string     `json:"-" gorm:"not null"`
	FirstName    string     `json:"first_name" gorm:"not null"`
	LastName     string     `json:"last_name" gorm:"not null"`
	Role         string     `json:"role" gorm:"not null;default:'employee'"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Client specific fields
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`

	// Relations
	Cars         []Car         `json:"cars,omitempty"`
	Appointments []Appointment `json:"appointments,omitempty"`
}

// TableName especifica o nome da tabela
func (User) TableName() string {
	return "users"
}

func NewUser(email, password, firstName, lastName, role string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}
	if firstName == "" {
		return nil, errors.New("first name is required")
	}
	if lastName == "" {
		return nil, errors.New("last name is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if role == "" {
		role = "employee"
	}

	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FirstName:    firstName,
		LastName:     lastName,
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// User roles
const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleEmployee = "employee"
	RoleClient   = "client" // New role
)

// ValidateRole checks if the role is valid
func (u *User) ValidateRole() bool {
	switch u.Role {
	case RoleAdmin, RoleManager, RoleEmployee, RoleClient:
		return true
	default:
		return false
	}
}

// IsClient returns true if user is a client
func (u *User) IsClient() bool {
	return u.Role == RoleClient
}

// IsEmployee returns true if user is employee, manager or admin
func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee || u.Role == RoleManager || u.Role == RoleAdmin
}

// CanManageUsers returns true if user can manage other users
func (u *User) CanManageUsers() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}
