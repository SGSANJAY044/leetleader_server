package main

import (
	"log"
	"leetleader_server/internal/database"
	"leetleader_server/internal/config"
	"leetleader_server/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	database.ConnectDatabase()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set Gin mode based on environment
	gin.SetMode(cfg.Environment)

	// Create Gin router
	r := gin.New()

	// Add default middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup routes
	routes.AuthRoutes(r)

	// Start server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}