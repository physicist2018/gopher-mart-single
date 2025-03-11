package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/physicist2018/gopher-mart-single/internal/entities"
	repository "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/user"
)

// RegisterController is a struct that represents a controller for registration.
// It contains a pointer to a RegisterUseCase object.
type RegisterController struct {
	registerUseCase user.RegisterUseCase
}

// NewRegisterController is a function that creates a new RegisterController.
// It takes a pointer to a RegisterUseCase object as an argument and returns a pointer to a RegisterController.
func NewRegisterController(registerUseCase user.RegisterUseCase) *RegisterController {
	return &RegisterController{registerUseCase: registerUseCase}
}

// Register is a function that handles the registration request.
// It takes a ResponseWriter and a Request as arguments.
// It decodes the request body into a User object, registers the user using the RegisterUseCase,
// and sets a cookie with the token in the response.
// If there is an error during the process, it returns a 400 Bad Request or 401 Unauthorized error.
func (c *RegisterController) Register(w http.ResponseWriter, r *http.Request) {
	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := c.registerUseCase.Execute(user.Login, user.Password)
	if errors.Is(err, repository.ErrUserAlreadyExists) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	} else if errors.Is(err, repository.ErrInternalServerError) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
