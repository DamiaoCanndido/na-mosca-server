package footballApi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/driven/foot_api_repo/dto"
	"github.com/DamiaoCanndido/na-mosca-server/internal/domain"
	"github.com/google/uuid"
	cache "github.com/patrickmn/go-cache"
)

const (
	baseURL = "https://v3.football.api-sports.io"
)

type FootballAPI struct {
	client   *http.Client
	apiKey   string
	cache    *cache.Cache
	poolRepo domain.PoolRepository
}

func NewFootballAPI(poolRepo domain.PoolRepository) *FootballAPI {
	apiKey := os.Getenv("FOOTBALL_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_API_KEY not found in environment variables")
	}

	return &FootballAPI{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey:   apiKey,
		cache:    cache.New(12*time.Hour, 24*time.Hour), // Cache com expiração de 12 horas
		poolRepo: poolRepo,
	}
}

func (api *FootballAPI) makeRequest(endpoint string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", baseURL, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Adiciona headers necessários
	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")
	req.Header.Add("x-rapidapi-key", api.apiKey)

	// Adiciona parâmetros de query
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// AddGameToPool implements domain.PoolRepository.
func (api *FootballAPI) AddGameToPool(apiGameID int, ownerID, poolID uuid.UUID) (*domain.Game, error) {
	// Verificar se a pool existe e pertence ao usuário
	pool, err := api.poolRepo.GetByID(poolID)
	if err != nil {
		return nil, fmt.Errorf("pool not found: %v", err)
	}

	if pool.OwnerID != ownerID {
		return nil, fmt.Errorf("user is not the owner of this pool")
	}

	// Verifica no cache primeiro
	cacheKey := fmt.Sprintf("game:%d", apiGameID)
	if cachedGame, found := api.cache.Get(cacheKey); found {
		return cachedGame.(*domain.Game), nil
	}

	// Busca o jogo na API externa
	params := map[string]string{
		"id": fmt.Sprintf("%d", apiGameID),
	}

	resp, err := api.makeRequest("fixtures", params)
	if err != nil {
		return nil, fmt.Errorf("error fetching game with ID %d: %v", apiGameID, err)
	}
	defer resp.Body.Close()

	var apiResponse dto.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	responseSlice, ok := apiResponse.Response.([]any)
	if !ok || len(responseSlice) == 0 {
		return nil, fmt.Errorf("game not found with ID: %d", apiGameID)
	}

	var apiFixtures []dto.ApiFixture
	responseData, _ := json.Marshal(apiResponse.Response)
	if err := json.Unmarshal(responseData, &apiFixtures); err != nil {
		return nil, fmt.Errorf("error parsing fixture data: %v", err)
	}

	apiFixture := apiFixtures[0]
	game := &domain.Game{
		ID:        uuid.New(),
		ApiGameID: apiGameID,
		HomeTeam:  apiFixture.Teams.Home.Name,
		AwayTeam:  apiFixture.Teams.Away.Name,
		StartTime: func() time.Time {
			parsedDate, _ := time.Parse(time.RFC3339, apiFixture.Fixture.Date)
			return parsedDate
		}(),
		Status:    apiFixture.Fixture.Status.Short,
		HomeScore: apiFixture.Goals.Home,
		AwayScore: apiFixture.Goals.Away,
	}

	// Salvar o jogo no banco de dados
	if err := api.poolRepo.SaveGame(game); err != nil {
		return nil, fmt.Errorf("error saving game: %v", err)
	}

	// Associar o jogo à pool
	if err := api.poolRepo.AddGameToPool(game.ID, poolID); err != nil {
		return nil, fmt.Errorf("error associating game with pool: %v", err)
	}

	// Armazena no cache
	api.cache.Set(cacheKey, game, cache.DefaultExpiration)

	return game, nil
}

func (api *FootballAPI) GetLeagues(leagueIDs []int) ([]domain.League, error) {
	var leagues []domain.League

	for _, id := range leagueIDs {
		cacheKey := fmt.Sprintf("league:%d", id)
		if cachedData, found := api.cache.Get(cacheKey); found {
			leagues = append(leagues, cachedData.(domain.League))
			continue
		}

		params := map[string]string{
			"id": fmt.Sprintf("%d", id),
			"current": "true",
		}

		resp, err := api.makeRequest("leagues", params)
		if err != nil {
			return nil, fmt.Errorf("error fetching league with ID %d: %v", id, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response for league with ID %d: %v", id, err)
		}

		var apiResponse dto.APIResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return nil, fmt.Errorf("error decoding response for league with ID %d: %v", id, err)
		}

		if len(apiResponse.Errors) > 0 {
			return nil, fmt.Errorf("API errors for league with ID %d: %v", id, apiResponse.Errors)
		}

		var apiLeagues []dto.ApiLeague
		responseData, _ := json.Marshal(apiResponse.Response)
		if err := json.Unmarshal(responseData, &apiLeagues); err != nil {
			return nil, fmt.Errorf("erro ao decodificar liga com ID %d: %v", id, err)
		}

		for _, apiLeague := range apiLeagues {
			league := domain.League{
				ID:      apiLeague.League.ID,
				Name:    apiLeague.League.Name,
				Country: apiLeague.Country.Name,
				Season:  apiLeague.Seasons[0].Year,
				Code:    apiLeague.Country.Code,
				Logo:    apiLeague.League.Logo,
				Flag:    apiLeague.Country.Flag,
			}
			leagues = append(leagues, league)
			api.cache.Set(cacheKey, league, cache.DefaultExpiration)
		}
	}

	return leagues, nil
}

func (api *FootballAPI) GetFixturesByLeague(leagueID int, season string, status string) ([]domain.Fixture, error) {
	cacheKey := fmt.Sprintf("fixtures:%d:%s:%s", leagueID, season, status)
	if cachedData, found := api.cache.Get(cacheKey); found {
		return cachedData.([]domain.Fixture), nil
	}

	params := map[string]string{
		"league": fmt.Sprintf("%d", leagueID),
		"season": season,
		"status": status,
	}

	resp, err := api.makeRequest("fixtures", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse dto.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	if len(apiResponse.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", apiResponse.Errors)
	}

	var apiFixtures []dto.ApiFixture
	responseData, _ := json.Marshal(apiResponse.Response)
	if err := json.Unmarshal(responseData, &apiFixtures); err != nil {
		return nil, err
	}

	fixtures := make([]domain.Fixture, len(apiFixtures))
	for i, apiFixture := range apiFixtures {
		fixtures[i] = domain.Fixture{
			ID:   apiFixture.Fixture.ID,
			Date: func() time.Time {
				parsedDate, _ := time.Parse(time.RFC3339, apiFixture.Fixture.Date)
				return parsedDate
			}(),
			Status: apiFixture.Fixture.Status.Short,
			Season: apiFixture.League.Season,
			Round:  apiFixture.League.Round,
			HomeTeam: domain.Team{
				ID:   apiFixture.Teams.Home.ID,
				Name: apiFixture.Teams.Home.Name,
				Logo: apiFixture.Teams.Home.Logo,
			},
			AwayTeam: domain.Team{
				ID:   apiFixture.Teams.Away.ID,
				Name: apiFixture.Teams.Away.Name,
				Logo: apiFixture.Teams.Away.Logo,
			},
			GoalsHome: apiFixture.Goals.Home,
			GoalsAway: apiFixture.Goals.Away,
		}
	}

	api.cache.Set(cacheKey, fixtures, cache.DefaultExpiration)
	return fixtures, nil
}

func (api *FootballAPI) GetTodayFixtures() ([]domain.Fixture, error) {
	today := time.Now().Format("2006-01-02") // Format today's date as YYYY-MM-DD
	cacheKey := fmt.Sprintf("fixtures:today:%s", today)

	// Check cache
	if cachedData, found := api.cache.Get(cacheKey); found {
		return cachedData.([]domain.Fixture), nil
	}

	query_strings := map[string]string{"date": today, "status": "NS"}

	resp, err := api.makeRequest("fixtures", query_strings)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse dto.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	if len(apiResponse.Errors) > 0 {
		return nil, fmt.Errorf("API error: %v", apiResponse.Errors)
	}

	var apiFixtures []dto.ApiFixture
	responseData, _ := json.Marshal(apiResponse.Response)
	if err := json.Unmarshal(responseData, &apiFixtures); err != nil {
		return nil, err
	}

	fixtures := make([]domain.Fixture, len(apiFixtures))
	for i, apiFixture := range apiFixtures {
		fixtures[i] = domain.Fixture{
			ID:   apiFixture.Fixture.ID,
			Date: func() time.Time {
				parsedDate, _ := time.Parse(time.RFC3339, apiFixture.Fixture.Date)
				return parsedDate
			}(),
			Status: apiFixture.Fixture.Status.Short,
			HomeTeam: domain.Team{
				ID:   apiFixture.Teams.Home.ID,
				Name: apiFixture.Teams.Home.Name,
				Logo: apiFixture.Teams.Home.Logo,
			},
			AwayTeam: domain.Team{
				ID:   apiFixture.Teams.Away.ID,
				Name: apiFixture.Teams.Away.Name,
				Logo: apiFixture.Teams.Away.Logo,
			},
			GoalsHome: apiFixture.Goals.Home,
			GoalsAway: apiFixture.Goals.Away,
		}
	}

	api.cache.Set(cacheKey, fixtures, cache.DefaultExpiration)
	return fixtures, nil
}