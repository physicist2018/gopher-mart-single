package entities

import "time"

type Transaction struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	OrderID         string    `db:"order_id" json:"order_id"`
	TypeTransaction string    `db:"type_transaction" json:"type_transaction"`
	Amount          float64   `db:"amount" json:"amount"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}
