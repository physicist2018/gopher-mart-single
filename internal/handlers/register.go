package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
	"github.com/physicist2018/gopher-mart-single/internal/ports/repository"
)

type Handler struct {
	userRepo    repository.UserRepository
	authService authservice.AuthService
}

func NewHandler(userRepo repository.UserRepository, authService authservice.AuthService) *Handler {
	return &Handler{userRepo: userRepo,
		authService: authService}
}

func (h *Handler) RegisterUser(c echo.Context) error {
	// Регистрация пользователя

	var userInput models.User
	if err := c.Bind(&userInput); err != nil {

		return err
	}

	_, err := h.authService.Register(c.Request().Context(), userInput.Login, userInput.Password)
	if err != nil {
		return err
	}

	// Автоматическая авторизация после успешной регистрации
	token, err := h.authService.Login(c.Request().Context(), userInput.Login, userInput.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			echo.Map{"error": "Failed to generate token"},
		)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User registered",
		"token":   token, // Возвращаем JWT токен
	})

}
