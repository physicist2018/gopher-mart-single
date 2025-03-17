package withdraw

import (
	"context"
	"errors"
	"time"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	balancerepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/balance"
	withdrawrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/withdraw"
)

type WithdrawUseCase struct {
	balanceRepository  balancerepo.Repository
	withdrawRepository withdrawrepo.Repository
}

func NewWithdrawUseCase(balanceRepository balancerepo.Repository,
	withdrawRepository withdrawrepo.Repository) *WithdrawUseCase {
	return &WithdrawUseCase{
		balanceRepository:  balanceRepository,
		withdrawRepository: withdrawRepository,
	}
}

func (w *WithdrawUseCase) Withdraw(userID int, orderID string, orderPrice float64) (*entities.WithdrawRequest, error) {
	// Делаем попытку запросить списание баллов на новый заказ.
	// проверяем, есть ли у пользователя баллы
	ctx := context.Background()
	ctxBalance, cancelBalance := context.WithTimeout(ctx, time.Second*3)
	defer cancelBalance()

	balance, err := w.balanceRepository.GetBalanceForUserID(ctxBalance, userID)
	if errors.Is(err, balancerepo.ErrBalanceNotFound) {
		return nil, balancerepo.ErrBalanceNotFound // эта ошибка говорит о том, что у пользователя нет баланса, что эквивалентно отсутствию самого пользователя
	}

	if balance.Current < orderPrice {
		return nil, balancerepo.ErrBalanceIsNotEnough // эта ошибка говорит о том, что у пользователя не хватает баланса на списание
	}

	ctxWithdraw, cancelWithdraw := context.WithTimeout(ctx, time.Second*3)
	defer cancelWithdraw()
	err = w.withdrawRepository.CreateWithdraw(ctxWithdraw, userID, orderID, orderPrice)

	return nil, nil
}

func (w *WithdrawUseCase) Withdrawals(userID int) ([]entities.WithdrawResponse, error) {
	return nil, nil
}
