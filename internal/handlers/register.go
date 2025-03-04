package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/repository"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userRepo    repository.UserRepository
	authService authservice.AuthService
}

func NewHandler(userRepo repository.UserRepository, authService authservice.AuthService) *Handler {
	return &Handler{userRepo: userRepo,
		authService: authService}
}

func (h *Handler) RegisterUser(c *gin.Context) {
	// Регистрация пользователя

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Проверяем уникальность логина
	existingUser, err := h.userRepo.GetUserByLogin(user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": http.StatusText(http.StatusConflict)})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	user.Password = string(hashedPassword)

	// Сохраняем пользователя
	if err := h.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User registered",
	})

}
