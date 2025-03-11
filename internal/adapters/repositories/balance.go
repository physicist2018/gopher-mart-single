package repositories

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories"
)

type BalanceRepositoryImpl struct {
	db *sqlx.DB
}

const balanceQueryString = `
SELECT 
    u.balance + 
    COALESCE(SUM(CASE WHEN t.type_transaction = 'ACCRUAL' THEN t.amount ELSE 0 END), 0) - 
    COALESCE(SUM(CASE WHEN t.type_transaction = 'WITHDRAWAL' THEN t.amount ELSE 0 END), 0) AS remaining_balance,
    COALESCE(SUM(CASE WHEN t.type_transaction = 'WITHDRAWAL' THEN t.amount ELSE 0 END), 0) AS total_withdrawal
FROM users u
LEFT JOIN Transactions t ON u.id = t.user_id
WHERE u.id = $1
GROUP BY u.id;`

func (b *BalanceRepositoryImpl) GetBalanceForUserID(ctx context.Context, userId int) (*entities.Balance, error) {
	balance := &entities.Balance{}
	err := b.db.GetContext(ctx, balance, balanceQueryString, userId)
	if err != nil {
		return nil, errors.Join(repository.ErrUserNotFound, err)
	}
	return balance, nil
}
