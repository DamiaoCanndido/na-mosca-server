package main

import (
	"log"
	"os"

	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/driven/footballApi"
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/driven/userRepo"
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/handlers"
	"github.com/DamiaoCanndido/na-mosca-server/internal/adapters/drivers/http/routes"
	"github.com/DamiaoCanndido/na-mosca-server/internal/domain"
	"github.com/DamiaoCanndido/na-mosca-server/internal/ports"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(&domain.User{})

	// Initialize repositories
	userRepo := userRepo.NewUserRepository(db)
	footballRepo := footballApi.NewFootballAPI()

	// Initialize services
	userService := ports.NewUserService(userRepo)
	footballService := ports.NewFootballService(footballRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	footballHandler := handlers.NewFootballHandler(footballService)

	// Initialize router
	r := gin.Default()

	// Setup routes
	routes.SetupAuthRoutes(r, userHandler)
	routes.SetupFootballRoutes(r, footballHandler)

	// Start server
	r.Run(":8080")
} 