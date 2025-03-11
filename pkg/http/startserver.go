package http

import (
	"net/http"
	"strconv"

	"github.com/physicist2018/gopher-mart-single/internal/adapters/http/controllers"
	"github.com/physicist2018/gopher-mart-single/internal/interfaces/services"
	"github.com/physicist2018/gopher-mart-single/internal/usecases/user"
	"github.com/physicist2018/gopher-mart-single/pkg/middlewares"
)

type Server struct {
	ListenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		ListenAddr: listenAddr,
	}
}

func (s *Server) Start(authUseCase *user.AuthUseCase, registerUseCase *user.RegisterUseCase, tokenService services.TokenService) {
	authController := controllers.NewAuthController(authUseCase)
	registerController := controllers.NewRegisterController(registerUseCase)

	authMiddleware := middlewares.JWTAuthMiddleware(tokenService)

	http.HandleFunc("POST /api/user/register", registerController.Register)
	http.HandleFunc("POST /api/user/login", authController.Login)
	protected := http.NewServeMux()

	protected.HandleFunc("GET /protected", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middlewares.UserIDKey{}).(int)
		if ok {
			w.Write([]byte("Authenticated user ID: " + strconv.Itoa(userID)))
		} else {
			w.Write([]byte("Not Authenticated"))
		}
		w.Write([]byte("pppp"))
	})

	http.Handle("/", authMiddleware(protected))

	http.ListenAndServe(s.ListenAddr, nil)
}
