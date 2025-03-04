package models

import "time"

type User struct {
	ID          uint         `gorm:"primaryKey"`
	Login       string       `gorm:"unique;not null"`
	Password    string       `gorm:"not null"`
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	Balance     Balance      `gorm:"foreignKey:UserID"`
	Orders      []Order      `gorm:"foreignKey:UserID"`
	Withdrawals []Withdrawal `gorm:"foreignKey:UserID"` // Один ко многим (1 пользователь -> много списаний)
}
