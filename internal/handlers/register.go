package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
	"github.com/physicist2018/gopher-mart-single/internal/ports/repository"
)

type Handler struct {
	userRepo    repository.UserRepository
	authService authservice.AuthService
}

func NewHandler(userRepo repository.UserRepository, authService authservice.AuthService) *Handler {
	return &Handler{userRepo: userRepo,
		authService: authService}
}

var JwtKey = []byte("very-secret-passkey")

// User struct to represent a user in the system
type User struct {
	Username string
	Password string
}

// Credentials struct for handling login
type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

// Claims for JWT payload
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Register a new user
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.authService.Register(r.Context(), creds.Username, creds.Password)
	log.Println(err)
	if err != nil {
		if err == authservice.ErrUserAlreadyExists {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
