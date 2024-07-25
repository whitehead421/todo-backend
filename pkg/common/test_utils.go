package common

import (
	"github.com/whitehead421/todo-backend/pkg/entities"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() *gorm.DB {
	// Create a new in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		zap.L().Error("Failed to connect to database", zap.Error(err))
		panic(err)
	}

	err = db.AutoMigrate(&entities.User{}, &entities.Todo{})
	if err != nil {
		zap.L().Error("Failed to migrate tables", zap.Error(err))
		panic(err)
	}

	return db
}
