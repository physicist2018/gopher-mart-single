package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
)

// var JwtKey = []byte("GORM perform write (create/update/delete) operations run inside a transaction to ensure data consistency, you can disable it during initialization if it is not required")

// type Claims struct {
// 	UserID uint `json:"user_id"`
// 	jwt.StandardClaims
// }

func JWTAuthMiddleware(authService authservice.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Извлекаем токен из заголовка
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен с помощью AuthService
		user, err := authService.ValidateToken(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Устанавливаем информацию о пользователе в контекст Gin
		c.Set("userID", user.ID)
		c.Set("user", user)

		// Передаем управление следующему обработчику
		c.Next()
	}
}
