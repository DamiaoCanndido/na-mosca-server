package dto

type ApiCountry struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

type ApiLeague struct {
	League struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Logo string `json:"logo"`
	} `json:"league"`
	Country ApiCountry `json:"country"`
	Seasons []struct {
		Year     int    `json:"year"`
		Start    string `json:"start"`
		End      string `json:"end"`
		Current  bool   `json:"current"`
		Coverage struct {
			Fixtures struct {
				Events             bool `json:"events"`
				Lineups            bool `json:"lineups"`
				StatisticsFixtures bool `json:"statistics_fixtures"`
				StatisticsPlayers  bool `json:"statistics_players"`
			} `json:"fixtures"`
			Standings   bool `json:"standings"`
			Players     bool `json:"players"`
			TopScorers  bool `json:"top_scorers"`
			TopAssists  bool `json:"top_assists"`
			TopCards    bool `json:"top_cards"`
			Injuries    bool `json:"injuries"`
			Predictions bool `json:"predictions"`
			Odds        bool `json:"odds"`
		} `json:"coverage"`
	} `json:"seasons"`
}

type ApiFixture struct {
	Fixture struct {
		ID        int    `json:"id"`
		Referee   string `json:"referee"`
		Timezone  string `json:"timezone"`
		Date      string `json:"date"`
		Timestamp int64  `json:"timestamp"`
		Periods   struct {
			First  int64 `json:"first"`
			Second int64 `json:"second"`
		} `json:"periods"`
		Venue struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			City string `json:"city"`
		} `json:"venue"`
		Status struct {
			Long    string `json:"long"`
			Short   string `json:"short"`
			Elapsed int    `json:"elapsed"`
			Extra   int    `json:"extra"`
		} `json:"status"`
	} `json:"fixture"`
	League struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		Logo      string `json:"logo"`
		Flag      string `json:"flag"`
		Season    int    `json:"season"`
		Round     string `json:"round"`
		Standings bool   `json:"standings"`
	} `json:"league"`
	Teams struct {
		Home struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner *bool  `json:"winner"`
		} `json:"home"`
		Away struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Logo   string `json:"logo"`
			Winner *bool  `json:"winner"`
		} `json:"away"`
	} `json:"teams"`
	Goals struct {
		Home int `json:"home"`
		Away int `json:"away"`
	} `json:"goals"`
	Score struct {
		Halftime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"halftime"`
		Fulltime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"fulltime"`
		Extratime struct {
			Home *int `json:"home"`
			Away *int `json:"away"`
		} `json:"extratime"`
		Penalty struct {
			Home *int `json:"home"`
			Away *int `json:"away"`
		} `json:"penalty"`
	} `json:"score"`
}

type APIResponse struct {
	Get        string      `json:"get"`
	Parameters interface{} `json:"parameters"`
	Errors     []string    `json:"errors"`
	Results    int         `json:"results"`
	Response   interface{} `json:"response"`
}
