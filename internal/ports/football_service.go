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

func (s *FootballService) GetLeagues(country string) ([]domain.League, error) {
	return s.repo.GetLeagues(country)
}

func (s *FootballService) GetFixtures(leagueID int, season string) ([]domain.Fixture, error) {
	return s.repo.GetFixtures(leagueID, season)
}

func (s *FootballService) GetLiveFixtures() ([]domain.Fixture, error) {
	return s.repo.GetLiveFixtures()
} 