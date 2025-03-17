package entities

import "time"

type Order struct {
	ID        string    `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Status    string    `db:"status" json:"status"`
	Accrual   float64   `db:"accrual" json:"accrual,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
