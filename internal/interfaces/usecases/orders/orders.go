package orders

import "github.com/physicist2018/gopher-mart-single/internal/entities"

type UseCase interface {
	CreateOrder(userID int, orderID string) (*entities.Order, error)
	Orders(userID int) ([]entities.Order, error)
}
