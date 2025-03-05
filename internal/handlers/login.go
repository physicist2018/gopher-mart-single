package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
)

func (h *Handler) LoginUser(c echo.Context) error {
	// Аутентификация пользователя

	var creds models.User
	if err := c.Bind(&creds); err != nil {
		return err
	}

	token, err := h.authService.Login(c.Request().Context(), creds.Login, creds.Password)
	if err != nil {
		switch {
		case errors.Is(err, authservice.ErrUserNotFound), errors.Is(err, authservice.ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		default:
			return c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}

	// Возвращаем токен в ответе
	return c.JSON(http.StatusOK, gin.H{"token": token})
}
