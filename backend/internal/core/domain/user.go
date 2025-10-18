package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Password  string     `json:"-" gorm:"not null"`
	FirstName string     `gorm:"not null"`
	LastName  string     `gorm:"not null"`
	Role      string     `json:"role" gorm:"not null;default:'employee'"`
	IsActive  bool       `gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index"`
}

func NewUser(email, password, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
