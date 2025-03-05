package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/labstack/echo/v4"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/repository"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUser(t *testing.T) {
	// Настраиваем окружение

	r := echo.New()

	// Создаем тестовую базу данных
	db, err := setupDB("test.db")
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := authservice.NewAuthService("superkey", userRepo)
	h := NewHandler(userRepo, authService)

	r.POST("/api/user/login", h.LoginUser)

	// Добавляем пользователя в тестовую базу
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Login:    "testuser",
		Password: string(hashedPassword),
	}
	_ = userRepo.CreateUser(&user)

	// Тест на успешную аутентификацию
	body := map[string]string{
		"login":    "testuser",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Тест на неправильный пароль
	body["password"] = "wrongpassword"
	jsonBody, _ = json.Marshal(body)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Тест на неправильный логин
	body["login"] = "wronguser"
	body["password"] = "password123"
	jsonBody, _ = json.Marshal(body)
	req, _ = http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
