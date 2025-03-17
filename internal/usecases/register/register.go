package register

import (
	"context"
	"errors"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/user"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services/token"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUseCase struct {
	userRepo     repository.Repository
	tokenService token.TokenService
}

func NewRegisterUseCase(userRepo repository.Repository, tokenService token.TokenService) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (uc *RegisterUseCase) Execute(login, passHash string) (string, error) {
	user := &entities.User{
		Login:    login,
		Password: passHash,
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(passHash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashBytes)

	err = uc.userRepo.Save(context.Background(), user)

	if err != nil {
		return "", err
	}

	userRegidtered, err := uc.userRepo.FindByLogin(context.Background(), login)

	if err != nil {
		return "", err
	}

	token, err := uc.tokenService.GenerateToken(int(userRegidtered.ID))
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}
