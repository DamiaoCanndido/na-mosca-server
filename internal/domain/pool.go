package domain

import (
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name         string     `gorm:"not null"`
	OwnerID      uuid.UUID  `gorm:"not null"`
	Owner        User       `gorm:"foreignKey:OwnerID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Participants []User     `gorm:"many2many:pool_participants;"`
	Games		 []Game     `gorm:"many2many:pool_games;"`
}

type PoolRepository interface {
    Create(name string, ownerID uuid.UUID) (string, error)
	GetByID(poolID uuid.UUID) (*Pool, error)
    AddParticipant(poolID, userID uuid.UUID) (string, error)
	AddGame(gameID, userID uuid.UUID) (string, error)
}