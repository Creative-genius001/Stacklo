package db

import (
	"os"

	"github.com/Creative-genius001/Stacklo/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		logger.Fatal("DB_URL is not set in environment variables")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to the database:", err)
	}

	DB = db
}
