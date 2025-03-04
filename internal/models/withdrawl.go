package models

import "time"

type Withdrawal struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	OrderID   uint      `gorm:"not null"` // Внешний ключ на заказ
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
