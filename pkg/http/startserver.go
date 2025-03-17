package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services/order"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/physicist2018/gopher-mart-single/internal/adapters/http/controllers"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/auth"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/balance"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/orders"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/register"
	"github.com/physicist2018/gopher-mart-single/pkg/middlewares"
	"github.com/rs/zerolog"
)

// Server представляет HTTP-сервер с настройками маршрутов и зависимостями.
type Server struct {
	server          *http.Server
	authUseCase     auth.UseCase
	registerUseCase register.UseCase
	balanceUseCase  balance.UseCase
	orderUseCase    orders.UseCase
	tokenService    services.TokenService
	orderService    order.Service
	logger          *zerolog.Logger
}

// NewServer создает новый экземпляр HTTP-сервера с заданными параметрами.
// listenAddr - адрес, на котором будет запущен сервер.
// authUseCase - use case для аутентификации.
// regUseCase - use case для регистрации пользователей.
// balanceUseCase - use case для работы с балансом пользователя.
// orderUseCase - use case для работы с заказами.
// tokenService - сервис для работы с JWT-токенами.
// logger - логгер для записи событий.
// Возвращает указатель на созданный Server.
func NewServer(listenAddr string, authUseCase auth.UseCase,
	regUseCase register.UseCase, balanceUseCase balance.UseCase,
	orderUseCase orders.UseCase,
	tokenService services.TokenService,
	orderService order.Service,
	logger *zerolog.Logger) *Server {
	mux := http.NewServeMux()

	// Инициализация контроллеров
	authController := controllers.NewAuthController(authUseCase)
	registerController := controllers.NewRegisterController(regUseCase)
	balanceController := controllers.NewBalanceController(balanceUseCase)
	ordersController := controllers.NewOrdersController(orderUseCase, logger)

	// Middleware для аутентификации
	authMiddleware := middlewares.JWTAuthMiddleware(tokenService)

	// Публичные маршруты
	mux.HandleFunc("/api/user/register", registerController.Register)
	mux.HandleFunc("/api/user/login", authController.Login)

	// Защищенные маршруты
	mux.Handle("/api/user/balance", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			balanceController.Balance(w, r)
		}
	})))

	mux.Handle("/api/user/orders", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ordersController.Orders(w, r)
		case http.MethodPost:
			ordersController.CreateOrder(w, r)
		}
	})))

	// Оборачиваем mux в gzip middleware
	handler := middlewares.GzipMiddleware(mux)

	return &Server{
		server: &http.Server{
			Addr:    listenAddr,
			Handler: handler,
		},
		authUseCase:     authUseCase,
		registerUseCase: regUseCase,
		balanceUseCase:  balanceUseCase,
		orderUseCase:    orderUseCase,
		tokenService:    tokenService,
		orderService:    orderService,
		logger:          logger,
	}
}

// Start запускает HTTP-сервер и обрабатывает graceful shutdown при получении сигналов SIGINT или SIGTERM.
// Сервер запускается в отдельной горутине, а основная горутина ожидает сигналов для завершения работы.
func (s *Server) Start() {
	// Запуск сервера в отдельной горутине
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Err(fmt.Errorf("listen and serve: %w", err)).Msg("Server error")
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	// Обработка graceful shutdown при получении сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.orderService.CheckOrdersForAccrual()
				if err != nil {
					s.logger.Error().Err(err).Msg("Failed to check orders for accrual")
				}

			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit

	s.logger.Info().Msg("Shutting down server...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error().Err(fmt.Errorf("server shutdown: %w", err)).Msg("Server shutdown error")
	}

	s.logger.Info().Msg("Server exited")
}
