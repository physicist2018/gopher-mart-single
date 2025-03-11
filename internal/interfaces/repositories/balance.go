package repository

import (
	"context"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

type BalanceRepository interface {
	GetBalanceForUserID(ctx context.Context, userId int) (*entities.Balance, error)
}
