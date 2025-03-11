package balance

import (
	"context"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories"
)

type BalanceUseCase struct {
	balanceRepository repository.BalanceRepository
}

func NewBalanceUseCase(balanceRepo repository.BalanceRepository) *BalanceUseCase {
	return &BalanceUseCase{balanceRepository: balanceRepo}
}

func (b *BalanceUseCase) Execute(userID int) (*entities.Balance, error) {
	balance, err := b.balanceRepository.GetBalanceForUserID(context.Background(), userID)
	if err != nil {
		return nil, repository.ErrInternalServerError
	}

	return balance, nil
}
