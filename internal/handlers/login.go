package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	log.Println(err)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// // Check if the user exists
	// user, err := h.userRepo.GetUserByLogin(r.Context(), creds.Username)
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// // Compare the provided password with stored hashed password
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// // Create JWT token
	expirationTime := time.Now().Add(5 * time.Minute)
	// claims := &Claims{
	// 	Username: creds.Username,
	// 	StandardClaims: jwt.StandardClaims{
	// 		ExpiresAt: expirationTime.Unix(),
	// 	},
	// }

	// // Generate encoded token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// tokenString, err := token.SignedString(JwtKey)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// Set JWT as cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
	})

	// Set JWT in header
	w.Header().Set("Authorization", "Bearer "+token)

	w.Write([]byte("Login successful"))
}
