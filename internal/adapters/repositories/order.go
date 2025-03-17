package repositories

import (
	"context"
	"errors"

	orderrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/order"
	userrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/user"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

// OrderRepositoryImpl реализует интерфейс для работы с заказами в базе данных.
type OrderRepositoryImpl struct {
	db *sqlx.DB // Подключение к базе данных
}

// NewOrderRepositoryImpl создает новый экземпляр OrderRepositoryImpl.
// db - подключение к базе данных.
// Возвращает указатель на OrderRepositoryImpl.
func NewOrderRepositoryImpl(db *sqlx.DB) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{db: db}
}

// insertOrderQuery - SQL-запрос для вставки нового заказа в таблицу orders.
const insertOrderQuery string = `
INSERT INTO orders(id, user_id) VALUES(:id, :user_id)
`

// Save сохраняет заказ в базе данных.
// ctx - контекст для управления таймаутами и отменой.
// order - заказ, который нужно сохранить.
// Возвращает ошибку, если заказ уже существует или произошла внутренняя ошибка сервера.
func (o *OrderRepositoryImpl) Save(ctx context.Context, order *entities.Order) error {
	_, err := o.db.NamedExecContext(ctx, insertOrderQuery, order)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" { // 23505 - это код ошибки для уникальных ограничений (duplicate key)
			return orderrepo.ErrOrderAlreadyExists
		}
		return userrepo.ErrInternalServerError
	}
	return nil
}

// selectOrderQuery - SQL-запрос для получения заказа по его ID.
const selectOrderQuery string = `
SELECT * FROM orders WHERE id = $1 LIMIT 1;
`

// GetOrderByID возвращает заказ по его ID.
// ctx - контекст для управления таймаутами и отменой.
// orderID - ID заказа, который нужно найти.
// Возвращает заказ или ошибку, если заказ не найден.
func (o *OrderRepositoryImpl) GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error) {
	order := &entities.Order{}
	err := o.db.GetContext(ctx, order, selectOrderQuery, orderID)
	if err != nil { // order not found
		return nil, orderrepo.ErrOrderNotFound
	}
	return order, nil
}

// selectOrdersByUserIDQuery - SQL-запрос для получения всех заказов пользователя по его ID.
const selectOrdersByUserIDQuery string = `
SELECT * FROM orders WHERE user_id = $1 ORDER by created_at DESC;
`

// GetAllByUserID возвращает все заказы пользователя по его ID.
// ctx - контекст для управления таймаутами и отменой.
// userID - ID пользователя, чьи заказы нужно найти.
// Возвращает список заказов или ошибку, если произошла внутренняя ошибка сервера.
func (o *OrderRepositoryImpl) GetAllByUserID(ctx context.Context, userID int) ([]entities.Order, error) {
	var orders []entities.Order
	err := o.db.SelectContext(ctx, &orders, selectOrdersByUserIDQuery, userID)
	if err != nil {
		return nil, userrepo.ErrInternalServerError
	}
	return orders, nil
}

const updateOrderStatusQuery string = `
UPDATE orders SET status = $1, accrual = $2 WHERE id = $3;
`

func (o *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderID string, orderStatus string, orderAccrual float64) error {
	_, err := o.db.ExecContext(ctx, updateOrderStatusQuery, orderStatus, orderAccrual, orderID)
	if err != nil {
		return userrepo.ErrInternalServerError
	}
	return nil
}

const createAccrualQuery string = `
INSERT INTO transactions(user_id, order_id, type_transaction, amount) VALUES ($1, $2, $3, $4);
`

func (o *OrderRepositoryImpl) CreateAccrual(ctx context.Context, userID string, orderID string, amount float64) error {
	_, err := o.db.ExecContext(ctx, createAccrualQuery, userID, orderID, "ACCRUAL", amount)
	if err != nil {
		return userrepo.ErrInternalServerError
	}
	return nil
}

const getPendingOrders string = `
SELECT * FROM orders WHERE status in ('NEW', 'IN_PROGRESS');
`

func (o *OrderRepositoryImpl) GetPendingOrders(ctx context.Context) ([]entities.Order, error) {
	var pendingOrders []entities.Order
	err := o.db.SelectContext(ctx, &pendingOrders, getPendingOrders)
	if err != nil {
		return nil, userrepo.ErrInternalServerError
	}
	return pendingOrders, err
}
