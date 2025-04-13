package domain

import (
	"time"

	"github.com/google/uuid"
)

type PlanType string

const (
	PlanFree    PlanType = "FREE"
	PlanPremium PlanType = "PREMIUM"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	
	Plan           PlanType   `gorm:"type:varchar(10);default:FREE" json:"plan"`
	PlanExpiresAt  *time.Time `json:"plan_expires_at,omitempty"`

	TotalPoints    int        `gorm:"default:0" json:"total_points"`
	Credits        int        `gorm:"default:0" json:"credits"`

	Level          int        `gorm:"default:1" json:"level"`
	AvatarUrl      *string    `json:"avatar_url,omitempty"`

	IsActive       bool       `gorm:"default:true" json:"is_active"`
}

type UserRepository interface {
	RegisterUser(user *User) error
	FindByEmail(email string) (*User, error)
	VerifyPassword(user *User, password string) error
}

type UserService interface {
	RegisterUser(name, avatar_url, email, password string) (*User, error)
	Authenticate(email, password string) (string, error)
} 