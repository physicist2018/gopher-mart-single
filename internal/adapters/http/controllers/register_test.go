package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegisterUseCase is a mock implementation of the RegisterUseCase interface.
type MockRegisterUseCase struct {
	mock.Mock
}

func (m *MockRegisterUseCase) Execute(login, password string) (string, error) {
	args := m.Called(login, password)
	return args.String(0), args.Error(1)
}

func TestRegisterController_Register_Success(t *testing.T) {
	// Create a mock RegisterUseCase
	mockRegisterUseCase := new(MockRegisterUseCase)

	// Set up the expected behavior for the mock
	mockRegisterUseCase.On("Execute", "testuser", "testpassword").Return("mock-token", nil)

	// Create the RegisterController with the mock RegisterUseCase
	registerController := NewRegisterController(mockRegisterUseCase)

	// Create a test user
	user := entities.User{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Marshal the user into JSON
	userJSON, _ := json.Marshal(user)

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Register method
	registerController.Register(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	// Assert the response headers
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Expected Content-Type header to be application/json")
	assert.Equal(t, "Bearer mock-token", w.Header().Get("Authorization"), "Expected Authorization header to contain the token")

	// Assert the cookie was set
	cookies := w.Result().Cookies()
	defer w.Result().Body.Close()

	assert.Equal(t, 1, len(cookies), "Expected one cookie to be set")
	assert.Equal(t, "token", cookies[0].Name, "Expected cookie name to be 'token'")
	assert.Equal(t, "mock-token", cookies[0].Value, "Expected cookie value to be 'mock-token'")

	// Verify that the mock was called as expected
	mockRegisterUseCase.AssertExpectations(t)
}

func TestRegisterController_Register_InvalidRequestBody(t *testing.T) {
	// Create a mock RegisterUseCase
	mockRegisterUseCase := new(MockRegisterUseCase)

	// Create the RegisterController with the mock RegisterUseCase
	registerController := NewRegisterController(mockRegisterUseCase)

	// Create a test HTTP request with an invalid body
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Register method
	registerController.Register(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

	// Verify that the mock was not called
	mockRegisterUseCase.AssertNotCalled(t, "Execute")
}

func TestRegisterController_Register_UserAlreadyExists(t *testing.T) {
	// Create a mock RegisterUseCase
	mockRegisterUseCase := new(MockRegisterUseCase)

	// Set up the expected behavior for the mock
	mockRegisterUseCase.On("Execute", "existinguser", "testpassword").Return("", repository.ErrUserAlreadyExists)

	// Create the RegisterController with the mock RegisterUseCase
	registerController := NewRegisterController(mockRegisterUseCase)

	// Create a test user
	user := entities.User{
		Login:    "existinguser",
		Password: "testpassword",
	}

	// Marshal the user into JSON
	userJSON, _ := json.Marshal(user)

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Register method
	registerController.Register(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusConflict, w.Code, "Expected status code 409")

	// Assert the response body
	assert.Contains(t, w.Body.String(), repository.ErrUserAlreadyExists.Error(), "Expected error message in response body")

	// Verify that the mock was called as expected
	mockRegisterUseCase.AssertExpectations(t)
}

func TestRegisterController_Register_InternalServerError(t *testing.T) {
	// Create a mock RegisterUseCase
	mockRegisterUseCase := new(MockRegisterUseCase)

	// Set up the expected behavior for the mock
	mockRegisterUseCase.On("Execute", "testuser", "testpassword").Return("", repository.ErrInternalServerError)

	// Create the RegisterController with the mock RegisterUseCase
	registerController := NewRegisterController(mockRegisterUseCase)

	// Create a test user
	user := entities.User{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Marshal the user into JSON
	userJSON, _ := json.Marshal(user)

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Register method
	registerController.Register(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500")

	// Assert the response body
	assert.Contains(t, w.Body.String(), repository.ErrInternalServerError.Error(), "Expected error message in response body")

	// Verify that the mock was called as expected
	mockRegisterUseCase.AssertExpectations(t)
}
