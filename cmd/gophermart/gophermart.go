package main

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/physicist2018/gopher-mart-single/internal/adapters/repositories"
	"github.com/physicist2018/gopher-mart-single/internal/adapters/services"
	"github.com/physicist2018/gopher-mart-single/internal/config"
	"github.com/physicist2018/gopher-mart-single/internal/usecases/user"
	"github.com/physicist2018/gopher-mart-single/pkg/http"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting server ...")
	cfg := config.NewConfiguration()
	cfg.Parse()
	fmt.Println(cfg.String())
	server := http.NewServer(":8080")
	db, err := sqlx.Connect("postgres", "postgres://gen_user:CBukGT1984!!!@37.252.20.250:5432/default_db?sslmode=disable")
	if err != nil {
		log.Fatal().Err(err)
	}

	userRepo := repositories.NewUserRepositoryImpl(db)
	tokenService := services.NewTokenService("very-secret-key", time.Second*10)
	authUserCase := user.NewAuthUseCase(userRepo, tokenService)
	registerUseCase := user.NewRegisterUseCase(userRepo, tokenService)

	server.Start(authUserCase, registerUseCase, tokenService)
}
