package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/physicist2018/gopher-mart-single/internal/handlers"
)

type UserName struct{}

// Middleware to validate JWT token from header or cookie
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		// Check token in Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Check token in cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				if err == http.ErrNoCookie {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tokenStr = cookie.Value
		}

		// Initialize Claims
		claims := &handlers.Claims{}

		// Parse the JWT
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Set user information in request context
		r = r.WithContext(contextWithUser(r.Context(), claims.Username))

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// Context helpers for user information
func contextWithUser(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, UserName{}, username)
}

// func userFromContext(ctx context.Context) string {
// 	username, _ := ctx.Value(UserName{}).(string)
// 	return username
// }
