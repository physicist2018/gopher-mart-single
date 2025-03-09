package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/physicist2018/gopher-mart-single/internal/config"
	db "github.com/physicist2018/gopher-mart-single/internal/database/db/postgres"
	"github.com/physicist2018/gopher-mart-single/internal/handlers"
	"github.com/physicist2018/gopher-mart-single/internal/repository"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
)

func main() {
	// Миграция модели User

	cfg := config.LoadConfig()

	dbase := db.NewDB(cfg.DatabaseURI)
	defer dbase.Close()

	// queries := db.New(dbase)

	userRepo := repository.NewUserRepository(dbase)
	authService := authservice.NewAuthService(cfg.JWTSecret, userRepo)

	//	authMiddleware := middlewares.JWTAuthMiddleware(authService)

	handlers := handlers.NewHandler(userRepo, authService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", handlers.Register)
	mux.HandleFunc("POST /api/user/login", handlers.Login)

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
