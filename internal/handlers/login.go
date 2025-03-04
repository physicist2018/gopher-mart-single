package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
)

func (h *Handler) LoginUser(c *gin.Context) {
	// Аутентификация пользователя

	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), creds.Login, creds.Password)
	if err != nil {
		switch {
		case errors.Is(err, authservice.ErrUserNotFound), errors.Is(err, authservice.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Возвращаем токен в ответе
	c.JSON(http.StatusOK, gin.H{"token": token})
}
