package repositories

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) FindByLogin(ctx context.Context, login string) (*entities.User, error) {
	user := &entities.User{}
	err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		return nil, errors.Join(repository.ErrUserNotFound, err)
	}
	return user, nil
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *entities.User) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO users (login, password) VALUES (:login, :password)", user)

	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" { // 23505 - это код ошибки для уникальных ограничений (duplicate key)
			return errors.Join(repository.ErrUserAlreadyExists, err)
		}
		return repository.ErrInternalServerError
	}
	return err
}
