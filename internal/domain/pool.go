package domain

import (
	"time"

	"github.com/google/uuid"
)

type Pool struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name         string     `json:"name" gorm:"not null"`
	OwnerID      uuid.UUID  `json:"owner_id" gorm:"type:uuid;not null"`
	Owner        User       `json:"owner" gorm:"foreignKey:OwnerID"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Participants []User     `gorm:"many2many:pool_participants;"`
	Games		 []Game     `gorm:"many2many:pool_games;"`
}

type PoolRepository interface {
    Create(name string, ownerID uuid.UUID) (string, error)
	GetByID(poolID uuid.UUID) (*Pool, error)
    AddParticipant(poolID, userID uuid.UUID) (string, error)
	SaveGame(game *Game) error
	AddGameToPool(gameID, poolID uuid.UUID) error
}