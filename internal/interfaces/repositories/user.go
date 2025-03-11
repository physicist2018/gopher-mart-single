package repository

import (
	"context"
	"errors"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserAlreadyExists   = errors.New("user alreadey exists")
	ErrInternalServerError = errors.New("internal server error")
)

type UserRepository interface {
	FindByLogin(ctx context.Context, login string) (*entities.User, error)
	Save(ctx context.Context, user *entities.User) error
}
