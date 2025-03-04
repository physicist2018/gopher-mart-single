package models

import "time"

type Transaction struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Type      string    `gorm:"not null"`
	OrderID   *uint     `gorm:"foreignKey:OrderID"` // Внешний ключ на заказ (опционально)
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
