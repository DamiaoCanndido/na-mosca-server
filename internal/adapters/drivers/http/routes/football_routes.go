package routes

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/handlers"
	"github.com/gin-gonic/gin"
)

func SetupFootballRoutes(router *gin.Engine, footballHandler *handlers.FootballHandler) {
	// Grupo de rotas de futebol
	football := router.Group("/football")
	{
		// Rotas públicas
		football.GET("/leagues", footballHandler.GetLeagues)
		football.GET("/leagues/:leagueID/fixtures", footballHandler.GetFixtures)
		football.GET("/today", footballHandler.GetTodayFixtures)
	}
}