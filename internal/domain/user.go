package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRepository interface {
	RegisterUser(user *User) error
	FindByEmail(email string) (*User, error)
	VerifyPassword(user *User, password string) error
}

type UserService interface {
	RegisterUser(name, email, password string) (*User, error)
	Authenticate(email, password string) (string, error)
} 