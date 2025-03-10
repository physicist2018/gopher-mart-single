package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
)

// Authenticate user and set JWT in cookie and header
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), creds.Username, creds.Password)

	if err != nil {
		if errors.Is(err, authservice.ErrInvalidCredentials) {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(5 * time.Minute)

	// Set JWT as cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	// Set JWT in header
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
