package common

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dsn string) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Fatal(
			"Failed to connect to database",
			zap.Error(err),
		)
	}

	zap.L().Info("Connected to database")
	DB = db
}
