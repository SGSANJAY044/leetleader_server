package database

import (
	"log"

	"leetleader_server/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Replace with your actual PostgreSQL DSN
	dsn := "host=localhost user=leetleader password=leetleaderpassword dbname=leetleader port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Automigrate models
	err = db.AutoMigrate(&models.Student{}, &models.Class{}, &models.Department{}, &models.Staff{}, &models.Question{}, &models.Assignment{}, &models.FriendsQuestions{})
	if err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	DB = db
}
