package routes

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/handlers"
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupPoolRoutes(router *gin.Engine, poolHandler *handlers.PoolHandler) {
	// Grupo de rotas de pools
	pools := router.Group("/pools")
	{
		// Rotas protegidas
		protected := pools.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/create", poolHandler.CreatePool)
		}
	}
} 