package poolRepo

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PoolRepository struct {
	db *gorm.DB
}

func NewPoolRepository(db *gorm.DB) *PoolRepository {
	return &PoolRepository{
		db: db,
	}
}

// AddParticipant implements domain.PoolRepository.
func (r *PoolRepository) AddParticipant(poolID uuid.UUID, userID uuid.UUID) (string, error) {
	panic("unimplemented")
}

// AddGame implements domain.PoolRepository.
func (r *PoolRepository) AddGame(gameID uuid.UUID, userID uuid.UUID) (string, error) {
	panic("unimplemented")
}

// GetByID implements domain.PoolRepository.
func (r *PoolRepository) GetByID(poolID uuid.UUID) (*domain.Pool, error) {
	var pool domain.Pool
	if err := r.db.First(&pool, "id = ?", poolID).Error; err != nil {
		return nil, err
	}
	return &pool, nil
}

func (r *PoolRepository) Create(name string, ownerID uuid.UUID) (string, error) {
	pool := &domain.Pool{
		ID:      uuid.New(),
		Name:    name,
		OwnerID: ownerID,
	}

	if err := r.db.Create(pool).Error; err != nil {
		return "", err
	}

	return pool.Name, nil
}
