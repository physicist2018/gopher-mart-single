package repository

import (
	"context"

	db "github.com/physicist2018/gopher-mart-single/internal/database/db/postgres"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id int32) (db.User, error)
	GetUserByLogin(ctx context.Context, login string) (db.User, error)
	UpdateUserBalance(ctx context.Context, arg db.UpdateUserBalanceParams) (db.User, error)
	UpdateUserRole(ctx context.Context, arg db.UpdateUserRoleParams) (db.User, error)
}
