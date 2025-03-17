package balance

import (
	"context"
	"errors"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

var (
	ErrBalanceNotFound    = errors.New("balance not found")
	ErrBalanceIsNotEnough = errors.New("not enough balance")
)

type Repository interface {
	GetBalanceForUserID(ctx context.Context, userID int) (*entities.Balance, error)
}
