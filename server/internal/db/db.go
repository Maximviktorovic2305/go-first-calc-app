package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Инициализация базы данных
func initDB() error {
	dsn := "host=localhost user=postgres password=admin dbname=go-calc port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Автоматическая миграция таблицы
	return db.AutoMigrate(&Calculation{})
}