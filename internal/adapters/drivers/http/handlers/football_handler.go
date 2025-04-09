package handlers

import (
	"net/http"
	"strconv"

	"github.com/DamiaoCanndido/na-mosca-server/internal/ports"
	"github.com/gin-gonic/gin"
)

type FootballHandler struct {
	service *ports.FootballService
}

func NewFootballHandler(service *ports.FootballService) *FootballHandler {
	return &FootballHandler{service: service}
}

func (h *FootballHandler) GetLeagues(c *gin.Context) {
	country := c.DefaultQuery("country", "brazil")

	leagues, err := h.service.GetLeagues(country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar ligas",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, leagues)
}

func (h *FootballHandler) GetFixtures(c *gin.Context) {
	leagueID, err := strconv.Atoi(c.Param("leagueID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ID da liga inválido",
			"details": err.Error(),
		})
		return
	}

	season := c.DefaultQuery("season", "2024")

	fixtures, err := h.service.GetFixtures(leagueID, season)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar jogos",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fixtures)
}

func (h *FootballHandler) GetLiveFixtures(c *gin.Context) {
	fixtures, err := h.service.GetLiveFixtures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar jogos ao vivo",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, fixtures)
} 