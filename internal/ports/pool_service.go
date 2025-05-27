package ports

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/domain"
	"github.com/google/uuid"
)

type PoolService struct {
	repo domain.PoolRepository
}

func NewPoolService(repo domain.PoolRepository) *PoolService {
	return &PoolService{repo: repo}
}

func (s *PoolService) Create(name string, ownerID uuid.UUID) (string, error) {
	return s.repo.Create(name, ownerID)
}
