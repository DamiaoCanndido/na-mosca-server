package domain

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
    ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    ApiGameID   int       `json:"api_game_id"`
    HomeTeam    string    `json:"home_team"`
    AwayTeam    string    `json:"away_team"`
    StartTime   time.Time `json:"start_time"`
    Status      string    `json:"status"`
    HomeScore   *int      `json:"home_score"`
    AwayScore   *int      `json:"away_score"`
}
