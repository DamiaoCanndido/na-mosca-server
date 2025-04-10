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
	cache "github.com/patrickmn/go-cache"
)

const (
	baseURL = "https://v3.football.api-sports.io"
)

type FootballAPI struct {
	client *http.Client
	apiKey string
	cache  *cache.Cache
}

func NewFootballAPI() *FootballAPI {
	apiKey := os.Getenv("FOOTBALL_API_KEY")
	if apiKey == "" {
		log.Fatal("FOOTBALL_API_KEY não encontrada nas variáveis de ambiente")
	}

	return &FootballAPI{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey: apiKey,
		cache:  cache.New(12*time.Hour, 24*time.Hour), // Cache com expiração de 12 horas
	}
}

func (api *FootballAPI) makeRequest(endpoint string, params map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", baseURL, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %v", err)
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
			return nil, fmt.Errorf("erro ao buscar liga com ID %d: %v", id, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler resposta para liga com ID %d: %v", id, err)
		}

		var apiResponse dto.APIResponse
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return nil, fmt.Errorf("erro ao decodificar resposta para liga com ID %d: %v", id, err)
		}

		if len(apiResponse.Errors) > 0 {
			return nil, fmt.Errorf("erros da API para liga com ID %d: %v", id, apiResponse.Errors)
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

func (api *FootballAPI) GetFixtures(leagueID int, season string, status string) ([]domain.Fixture, error) {
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

func (api *FootballAPI) GetLiveFixtures() ([]domain.Fixture, error) {
	resp, err := api.makeRequest("fixtures", map[string]string{"live": "all"})
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

	return fixtures, nil
}