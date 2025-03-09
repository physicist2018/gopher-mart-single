package repository

import (
	"context"
	"log"

	"github.com/lib/pq"
	db "github.com/physicist2018/gopher-mart-single/internal/database/db/postgres"
	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
	"github.com/physicist2018/gopher-mart-single/internal/ports/repository"
)

type userRepository struct {
	query *db.Queries
}

// Функция для создания нового репозитория пользователей
func NewUserRepository(dbtx db.DBTX) repository.UserRepository {
	return &userRepository{
		query: db.New(dbtx),
	}
}

// Реализация метода GetUserByLogin интерфейса UserRepository
func (r *userRepository) GetUserByLogin(ctx context.Context, login string) (db.User, error) {
	return r.query.GetUserByLogin(ctx, login)
}

func (r *userRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	log.Println(arg)
	user, err := r.query.CreateUser(ctx, arg)
	if err != nil {

		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // unique_violation
			return db.User{}, authservice.ErrUserAlreadyExists
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	return r.query.GetUserByID(ctx, id)
}

func (r *userRepository) UpdateUserBalance(ctx context.Context, arg db.UpdateUserBalanceParams) (db.User, error) {
	return r.query.UpdateUserBalance(ctx, arg)
}

func (r *userRepository) UpdateUserRole(ctx context.Context, arg db.UpdateUserRoleParams) (db.User, error) {
	return r.query.UpdateUserRole(ctx, arg)
}
