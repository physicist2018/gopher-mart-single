package main

import (
	"fmt"
	"github.com/physicist2018/gopher-mart-single/pkg/orderservice"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/physicist2018/gopher-mart-single/internal/adapters/repositories"
	"github.com/physicist2018/gopher-mart-single/internal/adapters/services"
	"github.com/physicist2018/gopher-mart-single/internal/config"
	"github.com/physicist2018/gopher-mart-single/internal/usecases/auth"
	"github.com/physicist2018/gopher-mart-single/internal/usecases/balance"
	"github.com/physicist2018/gopher-mart-single/internal/usecases/orders"
	"github.com/physicist2018/gopher-mart-single/internal/usecases/register"
	"github.com/physicist2018/gopher-mart-single/pkg/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	logger.Info().Msg("Starting server ...")
	cfg := config.NewConfiguration()
	cfg.Parse()
	logger.Info().Msg(cfg.String())

	db, err := sqlx.Connect("postgres", cfg.DatabaseURI)
	if err != nil {
		log.Fatal().Err(fmt.Errorf("connection error [%w]", err))
	}

	userRepo := repositories.NewUserRepositoryImpl(db)
	tokenService := services.NewTokenService("very-secret-key", time.Second*10)

	authUserCase := auth.NewAuthUseCase(userRepo, tokenService)
	registerUseCase := register.NewRegisterUseCase(userRepo, tokenService)
	balanceRepository := repositories.NewBalanceRepositoryImpl(db)
	balanceUseCase := balance.NewBalanceUseCase(balanceRepository)
	orderRepository := repositories.NewOrderRepositoryImpl(db)
	orderUseCase := orders.NewOrdersUseCase(orderRepository)
	orderService := orderservice.NewOrderService(orderRepository, cfg.AccrualSystemAddress, &logger)

	// Создаем сервер
	server := http.NewServer(cfg.RunAddress, authUserCase, registerUseCase,
		balanceUseCase, orderUseCase, tokenService, orderService, &logger)
	server.Start()
}
