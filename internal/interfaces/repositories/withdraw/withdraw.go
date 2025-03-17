package withdraw

import (
	"context"
	"errors"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

var (
	ErrWithdrawAlreadyExists = errors.New("withdraw already exists")
)

type Repository interface {
	CreateWithdraw(ctx context.Context, userID int, orderID string, priceTotal float64) error
	GetWithdrawByOrderID(ctx context.Context, orderID string) (*entities.WithdrawResponse, error)
	GetWithdrawals(ctx context.Context, userID int) ([]entities.WithdrawResponse, error)
}
