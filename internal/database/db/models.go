// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
	"time"
)

type Order struct {
	ID         int32
	UserID     int32
	Number     string
	Status     string
	Accrual    sql.NullFloat64
	UploadedAt time.Time
}

type User struct {
	ID       int32
	Login    string
	Password string
}
