package balance

import "github.com/physicist2018/gopher-mart-single/internal/entities"

type BalanceUseCase interface {
	Execute(userID int) (*entities.Balance, error)
}
