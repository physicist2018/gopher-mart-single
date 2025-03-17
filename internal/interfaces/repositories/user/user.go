package user

import (
	"context"
	"errors"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInternalServerError = errors.New("internal server error")
)

type Repository interface {
	FindByLogin(ctx context.Context, login string) (*entities.User, error)
	Save(ctx context.Context, user *entities.User) error
}
