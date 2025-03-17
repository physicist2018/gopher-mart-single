package withdraw

import "github.com/physicist2018/gopher-mart-single/internal/entities"

type UseCase interface {
	Withdraw(userID int, orderID string, orderPrice float64) (*entities.WithdrawResponse, error)
	Withdrawals(userID int) ([]entities.WithdrawResponse, error)
}
