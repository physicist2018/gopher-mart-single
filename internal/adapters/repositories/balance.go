package repositories

import (
	"context"
	"errors"

	balancerepository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/balance"

	"github.com/jmoiron/sqlx"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

// BalanceRepositoryImpl реализует интерфейс для работы с балансом пользователя в базе данных.
type BalanceRepositoryImpl struct {
	db *sqlx.DB // Подключение к базе данных
}

// NewBalanceRepositoryImpl создает новый экземпляр BalanceRepositoryImpl.
// db - подключение к базе данных.
// Возвращает реализацию интерфейса BalanceRepository.
func NewBalanceRepositoryImpl(db *sqlx.DB) balancerepository.Repository {
	return &BalanceRepositoryImpl{
		db: db,
	}
}

// balanceQueryString - SQL-запрос для получения текущего баланса и суммы списаний пользователя.
// Запрос вычисляет текущий баланс как сумму начального баланса пользователя и всех начислений,
// за вычетом всех списаний. Также возвращает общую сумму списаний.
const balanceQueryString = `
SELECT 
    u.balance + 
    COALESCE(SUM(CASE WHEN t.type_transaction = 'ACCRUAL' THEN t.amount ELSE 0 END), 0) - 
    COALESCE(SUM(CASE WHEN t.type_transaction = 'WITHDRAWAL' THEN t.amount ELSE 0 END), 0) AS current,
    COALESCE(SUM(CASE WHEN t.type_transaction = 'WITHDRAWAL' THEN t.amount ELSE 0 END), 0) AS withdrawn
FROM users u
LEFT JOIN Transactions t ON u.id = t.user_id
WHERE u.id = $1
GROUP BY u.id;`

// GetBalanceForUserID возвращает баланс пользователя по его ID.
// ctx - контекст для управления таймаутами и отменой.
// userID - ID пользователя, чей баланс нужно получить.
// Возвращает структуру Balance, содержащую текущий баланс и сумму списаний,
// или ошибку, если баланс не найден (ErrBalanceNotFound).
func (b *BalanceRepositoryImpl) GetBalanceForUserID(ctx context.Context, userID int) (*entities.Balance, error) {
	balance := &entities.Balance{}
	err := b.db.GetContext(ctx, balance, balanceQueryString, userID)
	if err != nil {
		return nil, errors.Join(balancerepository.ErrBalanceNotFound, err)
	}
	return balance, nil
}
