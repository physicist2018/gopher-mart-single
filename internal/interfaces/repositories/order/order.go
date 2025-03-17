package order

import (
	"context"
	"errors"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

var (
	ErrOrderNotFound                       = errors.New("order not found")
	ErrOrderAlreadyExists                  = errors.New("order alreadey exists")
	ErrOrderAlreadyUploadedByCurrentUserID = errors.New("order already uploaded by current user")
	ErrOrderAlreadyUploadedByAnotherUserID = errors.New("order already uploaded by another user")
	ErrOrderSuccsessfulyUploaded           = errors.New("order successfully uploaded")
)

type Repository interface {
	Save(ctx context.Context, order *entities.Order) error
	GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error)
	GetAllByUserID(ctx context.Context, userID int) ([]entities.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, orderStatus string, orderAccrual float64) error
	GetPendingOrders(ctx context.Context) ([]entities.Order, error)
	CreateAccrual(ctx context.Context, userID string, orderID string, amount float64) error
}
