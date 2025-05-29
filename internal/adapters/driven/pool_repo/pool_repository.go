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

// GetByID implements domain.PoolRepository.
func (r *PoolRepository) GetByID(poolID uuid.UUID) (*domain.Pool, error) {
	var pool domain.Pool
	if err := r.db.First(&pool, "id = ?", poolID).Error; err != nil {
		return nil, err
	}
	return &domain.Pool{
		ID:           pool.ID,
		Name:         pool.Name,
		Owner:        pool.Owner,
		OwnerID:      pool.OwnerID,
		Participants: pool.Participants,
		Games:        pool.Games,
	}, nil
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

func (r *PoolRepository) SaveGame(game *domain.Game) error {
	return r.db.Create(game).Error
}

func (r *PoolRepository) AddGameToPool(gameID, poolID uuid.UUID) error {
	return r.db.Table("pool_games").Create(map[string]interface{}{
		"pool_id": poolID,
		"game_id": gameID,
	}).Error
}
