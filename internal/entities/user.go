package entities

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Login     string    `db:"login" json:"login"`
	Password  string    `db:"password" json:"password"`
	Role      string    `db:"role" json:"role"`
	Balance   float64   `db:"balance" json:"balance"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
