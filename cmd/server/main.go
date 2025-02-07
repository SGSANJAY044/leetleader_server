package main

import (
	"github.com/gin-contrib/cors"
	"time"
	"leetleader_server/internal/config"
	"leetleader_server/internal/database"
	"leetleader_server/internal/routes"
	"log"

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
	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins (change for security)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add default middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Setup routes
	routes.AuthRoutes(r)
	routes.StudentRoutes(r)
	routes.StaffRoutes(r)

	// Start server
	if err := r.Run(cfg.ServerAddress); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}