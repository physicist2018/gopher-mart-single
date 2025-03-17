package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/auth"
)

// AuthController is a struct that represents a controller for authentication.
// It contains a pointer to an AuthUseCase object.
type AuthController struct {
	authUseCase auth.UseCase
}

// NewAuthController is a function that creates a new AuthController.
// It takes a pointer to an AuthUseCase object as an argument and returns a pointer to an AuthController.
func NewAuthController(authUseCase auth.UseCase) *AuthController {
	return &AuthController{authUseCase: authUseCase}
}

// Login is a function that handles the login request.
// It takes a ResponseWriter and a Request as arguments.
// It decodes the request body into a User object, validates the user using the AuthUseCase,
// and sets a cookie with the token in the response.
// If there is an error during the process, it returns a 400 Bad Request or 401 Unauthorized error.
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := c.authUseCase.Execute(user.Login, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Устанавливаем бессрочную куку, валидность токена зашита в нем самом
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
