package routes

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/handlers"
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/middleware"
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

		protected := football.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/pool/:pool_id/game/:api_game_id", footballHandler.AddGameToPool)
		}
	}
}