package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUseCase is a mock implementation of the AuthUseCase interface.
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Execute(login, password string) (string, error) {
	args := m.Called(login, password)
	return args.String(0), args.Error(1)
}
func TestAuthController_Login_Success(t *testing.T) {
	// Create a mock AuthUseCase
	mockAuthUseCase := new(MockAuthUseCase)

	// Set up the expected behavior for the mock
	mockAuthUseCase.On("Execute", "testuser", "testpassword").Return("mock-token", nil)

	// Create the AuthController with the mock AuthUseCase
	authController := NewAuthController(mockAuthUseCase)

	// Create a test user
	user := entities.User{
		Login:    "testuser",
		Password: "testpassword",
	}

	// Marshal the user into JSON
	userJSON, _ := json.Marshal(user)

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Login method
	authController.Login(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code, "Expected status code 200")

	// Assert the response headers
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Expected Content-Type header to be application/json")
	assert.Equal(t, "Bearer mock-token", w.Header().Get("Authorization"), "Expected Authorization header to contain the token")

	// Get the HTTP response from the recorder
	resp := w.Result()
	defer resp.Body.Close() // Ensure the response body is closed

	// Assert the cookie was set
	cookies := resp.Cookies()

	assert.Equal(t, 1, len(cookies), "Expected one cookie to be set")
	assert.Equal(t, "token", cookies[0].Name, "Expected cookie name to be 'token'")
	assert.Equal(t, "mock-token", cookies[0].Value, "Expected cookie value to be 'mock-token'")

	// Verify that the mock was called as expected
	mockAuthUseCase.AssertExpectations(t)
}

func TestAuthController_Login_InvalidRequestBody(t *testing.T) {
	// Create a mock AuthUseCase
	mockAuthUseCase := new(MockAuthUseCase)

	// Create the AuthController with the mock AuthUseCase
	authController := NewAuthController(mockAuthUseCase)

	// Create a test HTTP request with an invalid body
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Login method
	authController.Login(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected status code 400")

	// Verify that the mock was not called
	mockAuthUseCase.AssertNotCalled(t, "Execute")
}

func TestAuthController_Login_Unauthorized(t *testing.T) {
	// Create a mock AuthUseCase
	mockAuthUseCase := new(MockAuthUseCase)

	// Set up the expected behavior for the mock
	mockAuthUseCase.On("Execute", "testuser", "wrongpassword").Return("", errors.New("invalid credentials"))

	// Create the AuthController with the mock AuthUseCase
	authController := NewAuthController(mockAuthUseCase)

	// Create a test user
	user := entities.User{
		Login:    "testuser",
		Password: "wrongpassword",
	}

	// Marshal the user into JSON
	userJSON, _ := json.Marshal(user)

	// Create a test HTTP request
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	// Call the Login method
	authController.Login(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected status code 401")

	// Read the response body once and store it in a variable
	responseBody := w.Body.String()

	// Assert the response body
	assert.Contains(t, responseBody, "invalid credentials", "Expected error message in response body")

	// Verify that the mock was called as expected
	mockAuthUseCase.AssertExpectations(t)
}
