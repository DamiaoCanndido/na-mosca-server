package handlers

import (
	"net/http"
	"strconv"

	"github.com/DamiaoCanndido/na-mosca-server/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FootballHandler struct {
	service *ports.FootballService
}

func NewFootballHandler(service *ports.FootballService) *FootballHandler {
	return &FootballHandler{service: service}
}

func (h *FootballHandler) GetLeagues(c *gin.Context) {
	

	var leagueIDs = []int{71, 72, 15}
	

	leagues, err := h.service.GetLeagues(leagueIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar ligas", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leagues)
}

func (h *FootballHandler) GetFixtures(c *gin.Context) {
	leagueID, err := strconv.Atoi(c.Param("leagueID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid league ID",
			"details": err.Error(),
		})
		return
	}

	seasonQuery := c.Query("season")
	if seasonQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Season not provided",
			"details": "Season must be provided in the query string",
		})
		return
	}

	season := c.DefaultQuery("season", seasonQuery)
	status := c.DefaultQuery("status", "NS")

	fixtures, err := h.service.GetFixtures(leagueID, season, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error fetching fixtures",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fixtures)
}

func (h *FootballHandler) GetTodayFixtures(c *gin.Context) {
	fixtures, err := h.service.GetTodayFixtures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar jogos ao vivo",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fixtures)
}

func (h *FootballHandler) AddGameToPool(c *gin.Context) {
	apiGameID, err := strconv.Atoi(c.Param("api_game_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid game ID",
			"details": err.Error(),
		})
		return
	}

	poolID, err := uuid.Parse(c.Param("pool_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid pool ID",
			"details": err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "User ID not found",
			"details": "User ID must be provided",
		})
		return
	}

	ownerUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	game, err := h.service.AddGameToPool(apiGameID, ownerUUID, poolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error adding game to pool",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, game)
}