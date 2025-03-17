package balance

import (
	"context"

	balancerepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/balance"
	userrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/user"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

type BalanceUseCase struct {
	balanceRepository balancerepo.Repository
}

func NewBalanceUseCase(balanceRepo balancerepo.Repository) *BalanceUseCase {
	return &BalanceUseCase{balanceRepository: balanceRepo}
}

func (b *BalanceUseCase) Execute(userID int) (*entities.Balance, error) {
	balance, err := b.balanceRepository.GetBalanceForUserID(context.Background(), userID)
	if err != nil {
		return nil, userrepo.ErrInternalServerError
	}

	return balance, nil
}
