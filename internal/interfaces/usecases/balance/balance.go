package balance

import "github.com/physicist2018/gopher-mart-single/internal/entities"

type UseCase interface {
	Execute(userID int) (*entities.Balance, error)
}
