package main

import (
	"log"
	"os"

	"bolao/internal/adapters/driven/userRepo"
	"bolao/internal/adapters/drivers/http/handlers"
	"bolao/internal/adapters/drivers/http/middleware"
	"bolao/internal/domain"
	"bolao/internal/ports"

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

	// Initialize services
	userService := ports.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Initialize router
	r := gin.Default()

	// Public routes
	r.POST("/users", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Add protected routes here
	}

	// Start server
	r.Run(":8080")
} 