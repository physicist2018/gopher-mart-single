package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
)

// JWTAuthMiddleware для Echo
func JWTAuthMiddleware(authService authservice.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Получаем заголовок Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header is required"})
			}

			// Извлекаем токен из заголовка
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Валидируем токен с помощью AuthService
			user, err := authService.ValidateToken(c.Request().Context(), tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
			}

			// Устанавливаем информацию о пользователе в контекст Echo
			c.Set("userID", user.ID)
			c.Set("user", user)

			// Передаем управление следующему обработчику
			return next(c)
		}
	}
}
