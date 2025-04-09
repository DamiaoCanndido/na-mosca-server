package domain

import "time"

type League struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Code    string `json:"code"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
}

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type Fixture struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	HomeTeam  Team      `json:"homeTeam"`
	AwayTeam  Team      `json:"awayTeam"`
	GoalsHome int       `json:"goalsHome"`
	GoalsAway int       `json:"goalsAway"`
}

type FootballRepository interface {
	GetLeagues(country string) ([]League, error)
	GetFixtures(leagueID int, season string) ([]Fixture, error)
	GetLiveFixtures() ([]Fixture, error)
}

type FootballService interface {
	GetLeagues(country string) ([]League, error)
	GetFixtures(leagueID int, season string) ([]Fixture, error)
	GetLiveFixtures() ([]Fixture, error)
} 