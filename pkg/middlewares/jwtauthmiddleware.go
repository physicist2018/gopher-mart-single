package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services"
)

// UserIDKey is a struct that represents a key for the user ID in the context.
type UserIDKey struct{}

// JWTAuthMiddleware is a function that creates a middleware for JWT authentication.
// It takes a TokenService as an argument and returns a function that takes an http.Handler and returns an http.Handler.
// The returned function checks if the request has an Authorization header or a token cookie,
// validates the token using the TokenService, and adds the user ID to the request context.
// If there is an error during the process, it returns a 401 Unauthorized error.
func JWTAuthMiddleware(tokenService services.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				tokenParts := strings.Split(authHeader, " ")
				if len(tokenParts) == 2 || tokenParts[0] == "Bearer" {
					token = tokenParts[1]
				}
			}

			if token == "" {
				cookie, err := r.Cookie("token")
				if err == nil {
					token = cookie.Value
				}
			}

			if token == "" {
				http.Error(w, "Authorization token is required", http.StatusUnauthorized)
				return
			}

			userID, err := tokenService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
