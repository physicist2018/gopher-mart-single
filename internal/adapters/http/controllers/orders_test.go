package controllers_test

import (
	"bytes"
	"context"
	"errors"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/order"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/physicist2018/gopher-mart-single/internal/adapters/http/controllers"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
	"github.com/physicist2018/gopher-mart-single/pkg/middlewares"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrdersUseCase - mock для OrdersUseCase
type MockOrdersUseCase struct {
	mock.Mock
}

func (m *MockOrdersUseCase) CreateOrder(userID int, orderNumber string) (*entities.Order, error) {
	args := m.Called(userID, orderNumber)
	return args.Get(0).(*entities.Order), args.Error(1)
}

func (m *MockOrdersUseCase) Orders(userID int) ([]entities.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Order), args.Error(1)
}

// Тест для CreateOrder
func TestCreateOrder(t *testing.T) {
	logger := zerolog.Nop()

	t.Run("valid request - order created", func(t *testing.T) {
		mockUseCase := new(MockOrdersUseCase)
		controller := controllers.NewOrdersController(mockUseCase, &logger)

		// Подготовка запроса
		req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("0552199341654035"))
		req.Header.Set("Content-Type", "text/plain")

		// Подделка контекста с userID
		ctx := context.WithValue(req.Context(), middlewares.UserIDKey{}, 1)
		req = req.WithContext(ctx)

		// Успешное создание заказа
		mockUseCase.On("CreateOrder", 1, "0552199341654035").Return(&entities.Order{
			ID:     "0552199341654035",
			UserID: 1,
		}, nil)

		// Эмуляция HTTP-ответа
		rr := httptest.NewRecorder()

		// Вызов контроллера
		controller.CreateOrder(rr, req)

		// Проверка ответа
		assert.Equal(t, http.StatusCreated, rr.Code)
		mockUseCase.AssertCalled(t, "CreateOrder", 1, "0552199341654035")
	})

	t.Run("invalid order number", func(t *testing.T) {
		mockUseCase := new(MockOrdersUseCase)
		controller := controllers.NewOrdersController(mockUseCase, &logger)

		// Подготовка запроса с неверным номером заказа
		req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "text/plain")

		// Подделка контекста с userID
		ctx := context.WithValue(req.Context(), middlewares.UserIDKey{}, 1)
		req = req.WithContext(ctx)

		// Эмуляция HTTP-ответа
		rr := httptest.NewRecorder()

		// Вызов контроллера
		controller.CreateOrder(rr, req)

		// Проверка ответа
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	})

	t.Run("order already uploaded by current user", func(t *testing.T) {
		mockUseCase := new(MockOrdersUseCase)
		controller := controllers.NewOrdersController(mockUseCase, &logger)

		// Подготовка запроса
		req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString("0552199341654035"))
		req.Header.Set("Content-Type", "text/plain")

		// Подделка контекста с userID
		ctx := context.WithValue(req.Context(), middlewares.UserIDKey{}, 1)
		req = req.WithContext(ctx)

		// Симуляция ошибки: заказ уже загружен текущим пользователем
		mockUseCase.On("CreateOrder", 1, "0552199341654035").Return(&entities.Order{}, repository.ErrOrderAlreadyUploadedByCurrentUserID)

		// Эмуляция HTTP-ответа
		rr := httptest.NewRecorder()

		// Вызов контроллера
		controller.CreateOrder(rr, req)

		// Проверка ответа
		assert.Equal(t, http.StatusOK, rr.Code)
		mockUseCase.AssertCalled(t, "CreateOrder", 1, "0552199341654035")
	})
}

// Тест для Orders
func TestOrders(t *testing.T) {
	logger := zerolog.Nop()

	t.Run("valid request - orders found", func(t *testing.T) {
		mockUseCase := new(MockOrdersUseCase)
		controller := controllers.NewOrdersController(mockUseCase, &logger)

		// Подготовка запроса
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)

		// Подделка контекста с userID
		ctx := context.WithValue(req.Context(), middlewares.UserIDKey{}, 1)
		req = req.WithContext(ctx)

		// Симуляция успешного получения заказов
		mockUseCase.On("Orders", 1).Return([]entities.Order{
			{ID: "12345678903"},
		}, nil)

		// Эмуляция HTTP-ответа
		rr := httptest.NewRecorder()

		// Вызов контроллера
		controller.Orders(rr, req)

		// Проверка ответа
		assert.Equal(t, http.StatusOK, rr.Code)
		mockUseCase.AssertCalled(t, "Orders", 1)
	})

	t.Run("no orders found", func(t *testing.T) {
		mockUseCase := new(MockOrdersUseCase)
		controller := controllers.NewOrdersController(mockUseCase, &logger)

		// Подготовка запроса
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)

		// Подделка контекста с userID
		ctx := context.WithValue(req.Context(), middlewares.UserIDKey{}, 1)
		req = req.WithContext(ctx)

		// Симуляция ситуации, когда заказы не найдены
		mockUseCase.On("Orders", 1).Return([]entities.Order{}, nil)

		// Эмуляция HTTP-ответа
		rr := httptest.NewRecorder()

		// Вызов контроллера
		controller.Orders(rr, req)

		// Проверка ответа
		assert.Equal(t, http.StatusNoContent, rr.Code)
		mockUseCase.AssertCalled(t, "Orders", 1)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUseCase := new(MockOrdersUseCase)
		controller := controllers.NewOrdersController(mockUseCase, &logger)

		// Подготовка запроса
		req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)

		// Подделка контекста с userID
		ctx := context.WithValue(req.Context(), middlewares.UserIDKey{}, 1)
		req = req.WithContext(ctx)

		// Симуляция ошибки при получении заказов
		mockUseCase.On("Orders", 1).Return([]entities.Order{}, errors.New("internal server error"))

		// Эмуляция HTTP-ответа
		rr := httptest.NewRecorder()

		// Вызов контроллера
		controller.Orders(rr, req)

		// Проверка ответа
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockUseCase.AssertCalled(t, "Orders", 1)
	})
}
