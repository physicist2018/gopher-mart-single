package entities

import "time"

type Transaction struct {
	ID              int32     `db:"id" json:"id"`
	UserID          int32     `db:"user_id" json:"user_id"`
	OrderID         int32     `db:"order_id" json:"order_id"`
	TypeTransaction string    `db:"type_transaction" json:"type_transaction"`
	Amount          float64   `db:"amount" json:"amount"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}
