package entities

import (
	"time"
)

type User struct {
	ID        uint64    `gorm:"column:id;primary_key;auto_increment"`
	Email     string    `gorm:"column:email"`
	Name      string    `gorm:"column:name"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}