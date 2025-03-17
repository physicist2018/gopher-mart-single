package middlewares

import (
	"context"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services/token"
	"net/http"
	"strings"
)

// UserIDKey используется для хранения идентификатора пользователя в контексте запроса.
type UserIDKey struct{}

// JWTAuthMiddleware создает middleware для аутентификации пользователя с использованием JWT-токена.
// tokenService - сервис для работы с JWT-токенами (валидация, извлечение данных).
// Возвращает middleware, который проверяет наличие и валидность токена, а также добавляет userID в контекст запроса.
func JWTAuthMiddleware(tokenService token.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Извлекаем токен из заголовка Authorization (формат: "Bearer <token>")
				tokenParts := strings.Split(authHeader, " ")
				if len(tokenParts) == 2 || tokenParts[0] == "Bearer" {
					token = tokenParts[1]
				}
			}

			// Если токен не найден в заголовке, проверяем куки
			if token == "" {
				cookie, err := r.Cookie("token")
				if err == nil {
					token = cookie.Value
				}
			}

			// Если токен отсутствует, возвращаем ошибку 401 Unauthorized
			if token == "" {
				http.Error(w, "Authorization token is required", http.StatusUnauthorized)
				return
			}

			// Валидируем токен и извлекаем userID
			userID, err := tokenService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем userID в контекст запроса
			ctx := context.WithValue(r.Context(), UserIDKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
