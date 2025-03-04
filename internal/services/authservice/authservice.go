package authservice

import (
	"context"
	"errors"

	"github.com/physicist2018/gopher-mart-single/internal/models"
)

var (
	// ErrUserNotFound возвращается, если пользователь не найден.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials возвращается, если предоставлены неверные учетные данные.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserAlreadyExists возвращается, если пользователь уже существует.
	ErrUserAlreadyExists = errors.New("user already exists")
)

type AuthService interface {
	// Register регистрирует нового пользователя.
	Register(ctx context.Context, login, password string) (*models.User, error)

	// Login выполняет аутентификацию пользователя и возвращает токен.
	Login(ctx context.Context, login, password string) (string, error)

	// ValidateToken проверяет токен и возвращает информацию о пользователе.
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}
