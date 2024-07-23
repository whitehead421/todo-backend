package common

import (
	"github.com/whitehead421/todo-backend/pkg/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupTestDB() *gorm.DB {
	// Create a new in-memory SQLite database
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Todo{})

	return db
}
