package db

import (
	"github.com/Creative-genius001/Stacklo/services/payment/config"
	"github.com/Creative-genius001/go-logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.Cfg.DBUrl

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to the database:", err)
	}

	DB = db
}
