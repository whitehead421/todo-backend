package entities

import (
	"time"
)

type Todo struct {
	ID          uint64    `gorm:"column:id;primary_key;auto_increment"`
	Status      string    `gorm:"column:status"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
