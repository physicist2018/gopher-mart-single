package entities

import "time"

type Order struct {
	ID        int32     `db:"id" json:"id"`
	UserID    int32     `db:"user_id" json:"user_id"`
	Status    string    `db:"status" json:"status"`
	Accrual   float64   `db:"accrual" json:"accrual"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
