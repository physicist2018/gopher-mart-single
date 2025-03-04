package models

import "time"

type Order struct {
	ID           uint          `gorm:"primaryKey"`
	Number       string        `gorm:"unique;not null"`
	UserID       uint          `gorm:"not null"` // Внешний ключ на пользователя
	Status       string        `gorm:"not null"`
	Accrual      float64       `gorm:"default:0.0"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	Transactions []Transaction `gorm:"foreignKey:OrderID"` // Один ко многим (1 заказ -> много транзакций)
	Withdrawal   Withdrawal    `gorm:"foreignKey:OrderID"`
}
