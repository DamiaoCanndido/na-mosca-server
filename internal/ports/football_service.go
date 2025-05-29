package ports

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/domain"
	"github.com/google/uuid"
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
	return s.repo.GetFixturesByLeague(leagueID, season, status)
}

func (s *FootballService) GetTodayFixtures() ([]domain.Fixture, error) {
	return s.repo.GetTodayFixtures()
}

func (s *FootballService) AddGameToPool(apiGameID int, ownerID, poolID uuid.UUID) (*domain.Game, error) {
	game, err := s.repo.AddGameToPool(apiGameID, ownerID, poolID)
	if err != nil {
		return nil, err
	}
	return game, nil
}