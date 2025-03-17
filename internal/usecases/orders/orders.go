package orders

import (
	"context"
	"errors"
	orderrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/order"
	"time"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
)

type UseCase struct {
	ordersRepository orderrepo.Repository
}

func NewOrdersUseCase(ordersRepository orderrepo.Repository) *UseCase {
	return &UseCase{ordersRepository: ordersRepository}
}

func (o *UseCase) CreateOrder(userID int, orderID string) (*entities.Order, error) {

	// Проверить, есть ли уже такой заказ у этого пользователя
	// Проверить, есть ли уже такой заказ у другого пользователя
	// Создать заказ

	order := &entities.Order{
		UserID: userID,
		ID:     orderID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	orderOld, err := o.ordersRepository.GetOrderByID(ctx, orderID)
	if err == nil { //found order with given ID
		if orderOld.UserID != userID {
			return orderOld, orderrepo.ErrOrderAlreadyUploadedByAnotherUserID
		}
		return orderOld, orderrepo.ErrOrderAlreadyUploadedByCurrentUserID
	}

	// Если мы пришли сюда то соответсвующей лшибки не найдено
	if errors.Is(err, orderrepo.ErrOrderNotFound) {
		if err := o.ordersRepository.Save(ctx, order); err != nil {
			return nil, err
		}

	} else {
		return nil, err
	}

	return order, nil
}

func (o *UseCase) Orders(userID int) ([]entities.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return o.ordersRepository.GetAllByUserID(ctx, userID)
}
