package http

import (
	"net/http"

	"github.com/physicist2018/gopher-mart-single/internal/adapters/http/controllers"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/balance"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/user"
	"github.com/physicist2018/gopher-mart-single/pkg/middlewares"
	"github.com/rs/zerolog"
)

type Server struct {
	ListenAddr      string
	authUseCase     user.AuthUseCase
	registerUseCase user.RegisterUseCase
	balanceUseCase  balance.BalanceUseCase
	tokenService    services.TokenService
	logger          *zerolog.Logger
}

func NewServer(listenAddr string, authUseCase user.AuthUseCase,
	regUseCase user.RegisterUseCase, balanceUseCase balance.BalanceUseCase,
	tokenService services.TokenService,
	logger *zerolog.Logger) *Server {
	return &Server{
		ListenAddr:      listenAddr,
		authUseCase:     authUseCase,
		registerUseCase: regUseCase,
		balanceUseCase:  balanceUseCase,
		tokenService:    tokenService,
		logger:          logger,
	}
}

func (s *Server) Start() {
	authController := controllers.NewAuthController(s.authUseCase)
	registerController := controllers.NewRegisterController(s.registerUseCase)
	balanceController := controllers.NewBalanceController(s.balanceUseCase)

	authMiddleware := middlewares.JWTAuthMiddleware(s.tokenService)

	http.HandleFunc("POST /api/user/register", registerController.Register)
	http.HandleFunc("POST /api/user/login", authController.Login)

	balanceHandler := http.NewServeMux()
	balanceHandler.HandleFunc("GET /", balanceController.Balance)

	http.Handle("/api/user/balance", authMiddleware(balanceHandler))
	http.ListenAndServe(s.ListenAddr, nil)
}
