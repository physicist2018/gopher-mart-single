package auth

import (
	"context"
	"errors"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/user"

	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo     repository.Repository
	tokenService services.TokenService
}

func NewAuthUseCase(userRepo repository.Repository, tokenService services.TokenService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (uc *AuthUseCase) Execute(login, password string) (string, error) {
	user, err := uc.userRepo.FindByLogin(context.Background(), login)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Генерируем токен
	token, err := uc.tokenService.GenerateToken(int(user.ID))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil

}
