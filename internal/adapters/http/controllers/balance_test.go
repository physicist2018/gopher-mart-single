package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	"github.com/physicist2018/gopher-mart-single/pkg/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBalanceUseCase struct {
	mock.Mock
}

func (m *MockBalanceUseCase) Execute(userID int) (*entities.Balance, error) {
	args := m.Called(userID)
	return args[0].(*entities.Balance), args.Error(1)
}

// TestBalanceController_Balance_Success - тест для успешного получения баланса
func TestBalanceController_Balance_Success(t *testing.T) {
	// Создаем мок BalanceUseCase
	mockBalanceUseCase := &MockBalanceUseCase{}

	// Set up the expected behavior for the mock
	mockBalanceUseCase.On("Execute", 1).Return(&entities.Balance{
		Current:   10,
		Withdrawn: 100,
	}, nil)

	// Создаем контроллер с моком
	controller := NewBalanceController(mockBalanceUseCase)

	// Создаем HTTP-запрос с userID в контексте
	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	userID := 1
	ctx := req.Context()
	ctx = context.WithValue(ctx, middlewares.UserIDKey{}, userID)
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder для записи ответа
	w := httptest.NewRecorder()

	// Вызываем метод Balance
	controller.Balance(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	// Проверяем заголовок Content-Type
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Expected Content-Type header to be application/json")

	// Проверяем тело ответа
	var responseBalance entities.Balance
	err := json.NewDecoder(w.Body).Decode(&responseBalance)
	assert.NoError(t, err, "Failed to decode response body")
	assert.Equal(t, 10.0, responseBalance.Current, "Expected current to be 10")
	assert.Equal(t, 100.0, responseBalance.Withdrawn, "Expected withdrawn to be 100")
}

// TestBalanceController_Balance_NoUserKey - ест для случая отсутствия userID в контексте
func TestBalanceController_Balance_NoUserKey(t *testing.T) {
	// Создаем мок BalanceUseCase
	mockBalanceUseCase := &MockBalanceUseCase{}

	// Set up the expected behavior for the mock
	mockBalanceUseCase.On("Execute", 1).Return(&entities.Balance{
		Current:   10,
		Withdrawn: 100,
	}, nil)

	// Создаем контроллер с моком
	controller := NewBalanceController(mockBalanceUseCase)

	// Создаем HTTP-запрос с userID в контексте
	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	//userID := 1
	ctx := req.Context()
	//ctx = context.WithValue(ctx, middlewares.UserIDKey{}, userID)
	req = req.WithContext(ctx)

	// Создаем ResponseRecorder для записи ответа
	w := httptest.NewRecorder()

	// Вызываем метод Balance
	controller.Balance(w, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500")

}
