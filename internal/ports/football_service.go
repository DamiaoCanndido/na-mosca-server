package ports

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/domain"
)

type FootballService struct {
	repo domain.FootballRepository
}

func NewFootballService(repo domain.FootballRepository) *FootballService {
	return &FootballService{repo: repo}
}

func (s *FootballService) GetLeagues(leagueIDs []int) ([]domain.League, error) {
	return s.repo.GetLeagues(leagueIDs)
}

func (s *FootballService) GetFixtures(leagueID int, season string, status string) ([]domain.Fixture, error) {
	return s.repo.GetFixtures(leagueID, season, status)
}

func (s *FootballService) GetLiveFixtures() ([]domain.Fixture, error) {
	return s.repo.GetLiveFixtures()
}