package controllers

import (
	"encoding/json"
	"errors"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/order"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/user"
	"io"
	"net/http"

	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/orders"
	"github.com/physicist2018/gopher-mart-single/pkg/luhn"
	"github.com/physicist2018/gopher-mart-single/pkg/middlewares"
	"github.com/rs/zerolog"
)

type OrdersController struct {
	ordersUseCase orders.UseCase
	logger        *zerolog.Logger
}

func NewOrdersController(ordersUseCase orders.UseCase, logger *zerolog.Logger) *OrdersController {
	return &OrdersController{ordersUseCase: ordersUseCase, logger: logger}
}

func (o *OrdersController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey{}).(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		o.logger.Debug().Msg("Wrong content-type")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	numberBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		o.logger.Debug().Msg("Error reading body")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	numberStr := string(numberBytes)

	// Проверяем валидноть номера заказа
	if !luhn.ValidateLuhnNumber(numberStr) {
		o.logger.Debug().Msg("Error validating ordernum")
		http.Error(w, "неверный формат или номер заказа", http.StatusUnprocessableEntity)
		return
	}

	_, err = o.ordersUseCase.CreateOrder(userID, numberStr)

	if err != nil {
		o.logger.Debug().Err(err).Int("userID", userID).Send()
		switch {
		case errors.Is(err, order.ErrOrderAlreadyUploadedByCurrentUserID):
			http.Error(w, err.Error(), http.StatusOK)
			return
		case errors.Is(err, order.ErrOrderAlreadyUploadedByAnotherUserID):
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case errors.Is(err, repository.ErrInternalServerError):
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Error(w, order.ErrOrderSuccsessfulyUploaded.Error(), http.StatusCreated)
}

func (o *OrdersController) Orders(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey{}).(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orders, err := o.ordersUseCase.Orders(userID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	res, err := json.MarshalIndent(orders, "", "  ")
	if err != nil {
		json.NewEncoder(w).Encode(orders)
	}
	w.Write(res)
}
