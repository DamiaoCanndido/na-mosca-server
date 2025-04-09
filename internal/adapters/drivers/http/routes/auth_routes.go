package routes

import (
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/handlers"
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	// Grupo de rotas de autenticação
	auth := router.Group("/auth")
	{
		// Rotas públicas
		auth.POST("/register", userHandler.RegisterUser)
		auth.POST("/login", userHandler.Login)

		// Rotas protegidas
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Aqui podem ser adicionadas rotas protegidas relacionadas à autenticação
			// Por exemplo: refresh token, logout, etc.
		}
	}
} 