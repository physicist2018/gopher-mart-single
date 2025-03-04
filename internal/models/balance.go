package models

import "time"

type Balance struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"unique;not null"` // Внешний ключ на пользователя
	CurrentAmount float64   `gorm:"default:0.0"`
	BlockedAmount float64   `gorm:"default:0.0"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
