package db

import (
	"time"
)

// User struct definition
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	Username  string `gorm:"unique;not null"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
