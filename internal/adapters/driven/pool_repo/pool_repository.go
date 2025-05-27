package poolRepo

import "gorm.io/gorm"

type PoolRepository struct {
	db *gorm.DB
}

func NewPoolRepository(db *gorm.DB) *PoolRepository {
	return &PoolRepository{
		db: db,
	}
}

// func (r *PoolRepository) Create(name string, ownerID string) (string, error) {
	
// }