package db

import (
	"log"
	"server/internal/calculationService"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Инициализация базы данных
func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=admin dbname=go-calc port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&calculationService.Calculation{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

	return db, nil         
}