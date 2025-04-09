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
)

const (
	baseURL = "https://v3.football.api-sports.io"
)

type FootballAPI struct {
	client *http.Client
	apiKey string
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

// GetLeagues retorna as ligas disponíveis
func (api *FootballAPI) GetLeagues(country string) ([]domain.League, error) {
	params := map[string]string{
		"country": country,
	}

	resp, err := api.makeRequest("leagues", params)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar ligas: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var apiResponse dto.APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	if len(apiResponse.Errors) > 0 {
		return nil, fmt.Errorf("erros da API: %v", apiResponse.Errors)
	}

	var apiLeagues []dto.ApiLeague
	responseData, _ := json.Marshal(apiResponse.Response)
	if err := json.Unmarshal(responseData, &apiLeagues); err != nil {
		return nil, fmt.Errorf("erro ao decodificar ligas: %v", err)
	}

	leagues := make([]domain.League, len(apiLeagues))
	for i, apiLeague := range apiLeagues {
		leagues[i] = domain.League{
			ID:      apiLeague.League.ID,
			Name:    apiLeague.League.Name,
			Country: apiLeague.Country.Name,
			Code:    apiLeague.Country.Code,
			Logo:    apiLeague.League.Logo,
			Flag:    apiLeague.Country.Flag,
		}
	}

	return leagues, nil
}

// GetFixtures retorna os jogos de uma liga específica
func (api *FootballAPI) GetFixtures(leagueID int, season string) ([]domain.Fixture, error) {
	params := map[string]string{
		"league": fmt.Sprintf("%d", leagueID),
		"season": season,
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

	return fixtures, nil
}

// GetLiveFixtures retorna os jogos ao vivo
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