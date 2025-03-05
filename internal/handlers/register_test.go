package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/labstack/echo/v4"
	"github.com/physicist2018/gopher-mart-single/internal/database/connector"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/repository"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
)

func TestRegisterUser(t *testing.T) {
	// Test code here

	r := echo.New()

	db, err := setupDB("test.db")
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := authservice.NewAuthService("superkey", userRepo)
	h := NewHandler(userRepo, authService)

	r.POST("/api/user/register", h.RegisterUser)

	// Тест на успешную регистрацию
	body := map[string]string{
		"login":    "testuser",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Запускаем тест
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	// Тест на регистрацию с уже существующим логином
	req, _ = http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func setupDB(dbname string) (*connector.Connector, error) {
	os.Remove(dbname)
	db, err := connector.NewDBConnector("sqlite", dbname)
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{}, &models.Balance{}, &models.Order{}, &models.Transaction{}, &models.Withdrawal{})
	return db, err
}
