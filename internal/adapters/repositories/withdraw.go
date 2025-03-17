package repositories

import (
	"context"
	"errors"

	"github.com/lib/pq"
	orderrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/order"
	userrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/user"
	withdrawrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/withdraw"

	"github.com/jmoiron/sqlx"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

// WithdrawRepositoryImpl реализует интерфейс для работы с операциями списания (withdrawals) в базе данных.
type WithdrawRepositoryImpl struct {
	db *sqlx.DB // Подключение к базе данных
}

// NewWithdrawRepositoryImpl создает новый экземпляр WithdrawRepositoryImpl.
// db - подключение к базе данных.
// Возвращает указатель на WithdrawRepositoryImpl.
func NewWithdrawRepositoryImpl(db *sqlx.DB) *WithdrawRepositoryImpl {
	return &WithdrawRepositoryImpl{db: db}
}

// createWithdrawalTransaction - SQL-запрос для создания транзакции списания.
const createWithdrawalTransaction string = `
INSERT INTO transactions(user_id, order_id, type_transaction, amount) VALUES ($1, $2, 'WITHDRAWAL', $3);
`

// CreateWithdraw создает транзакцию списания для указанного пользователя.
// ctx - контекст для управления таймаутами и отменой.
// userID - ID пользователя, для которого создается списание.
// orderID - ID заказа, связанного с списанием.
// priceTotal - сумма списания.
// Возвращает ошибку, если:
//   - транзакция с таким orderID уже существует (ErrWithdrawAlreadyExists),
//   - произошла внутренняя ошибка сервера (ErrInternalServerError).
func (w *WithdrawRepositoryImpl) CreateWithdraw(ctx context.Context, userID int, orderID string, priceTotal float64) error {
	_, err := w.db.ExecContext(ctx, createWithdrawalTransaction, userID, orderID, priceTotal)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) { // Проверяем, является ли ошибка PostgreSQL
			if pqErr.Code == "23505" { // 23505 - код ошибки для нарушения уникальности (duplicate key)
				return errors.Join(withdrawrepo.ErrWithdrawAlreadyExists, err)
			}
		}
		return userrepo.ErrInternalServerError
	}
	return nil
}

// getWithdrawalsQuery - SQL-запрос для получения всех транзакций списания для указанного пользователя.
const getWithdrawalsQuery string = `
SELECT * FROM transactions WHERE user_id = $1 AND type_transaction = 'WITHDRAWAL' ORDER BY created_at DESC;
`

// GetWithdrawals возвращает список всех списаний для указанного пользователя.
// ctx - контекст для управления таймаутами и отменой.
// userID - ID пользователя, для которого запрашиваются списания.
// Возвращает список списаний в формате WithdrawResponse или ошибку, если произошла внутренняя ошибка сервера (ErrInternalServerError).
func (w *WithdrawRepositoryImpl) GetWithdrawals(ctx context.Context, userID int) ([]entities.WithdrawResponse, error) {
	var transactions []entities.Transaction
	err := w.db.SelectContext(ctx, &transactions, getWithdrawalsQuery, userID)
	if err != nil {
		return nil, userrepo.ErrInternalServerError
	}

	withdrawals := make([]entities.WithdrawResponse, len(transactions))
	for i := range transactions {
		withdrawals[i].OrderID = transactions[i].OrderID
		withdrawals[i].Sum = transactions[i].Amount
		withdrawals[i].ProcessedAt = transactions[i].CreatedAt
	}

	return withdrawals, nil
}

const getWithdrawByOrderIDQuery string = `
SELECT * FROM transactions WHERE order_id = $1 AND type_transaction = 'WITHDRAWAL';
`

func (w *WithdrawRepositoryImpl) GetWithdrawByOrderID(ctx context.Context, orderID string) (*entities.WithdrawResponse, error) {
	var transaction entities.Transaction
	err := w.db.GetContext(ctx, &transaction, getWithdrawByOrderIDQuery, orderID)
	if err != nil {
		return nil, orderrepo.ErrOrderNotFound
	}
	withdraw := &entities.WithdrawResponse{
		WithdrawRequest: entities.WithdrawRequest{
			OrderID: transaction.OrderID,
			Sum:     transaction.Amount,
		},

		ProcessedAt: transaction.CreatedAt,
	}
	return withdraw, nil
}
