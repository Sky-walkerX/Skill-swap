// cmd/server/main.go
package main

import (
	"log"

	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/database"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/middleware"
	apirouter "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.ConfigurableCORS(middleware.CORSWithFrontendURL(cfg.FrontendURL)))
	router.Use(middleware.SecurityHeaders())

	// Setup health routes
	apirouter.SetupHealthRoutes(router, &cfg)

	// Setup API routes
	api := router.Group("/api/v1")
	apirouter.SetupRoutes(api, db, &cfg)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
