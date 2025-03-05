package main

import (
	"log"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
	"github.com/physicist2018/gopher-mart-single/internal/config"
	"github.com/physicist2018/gopher-mart-single/internal/database/connector"
	"github.com/physicist2018/gopher-mart-single/internal/handlers"
	"github.com/physicist2018/gopher-mart-single/internal/middlewares"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/repository"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
)

func main() {
	// Миграция модели User

	cfg := config.LoadConfig()
	db, err := connector.NewDBConnector(cfg.DBType, cfg.DatabaseURI)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.User{}, &models.Balance{}, &models.Order{}, &models.Transaction{}, &models.Withdrawal{})

	r := echo.New()
	r.HideBanner = true
	r.Use(mid.LoggerWithConfig(
		mid.LoggerConfig{
			Format: "{time=${time_rfc3339}, id=${id} method=${method}, uri=${uri}, status=${status}}\n",
		},
	))
	r.Use(mid.Recover())

	log.Println("Secret:", cfg.JWTSecret)
	userRepo := repository.NewUserRepository(db)
	authService := authservice.NewAuthService(cfg.JWTSecret, userRepo)
	authMiddleware := middlewares.JWTAuthMiddleware(authService)

	handlers := handlers.NewHandler(userRepo, authService)

	r.POST("/api/user/register", handlers.RegisterUser)
	r.POST("/api/user/login", handlers.LoginUser)

	api := r.Group("/api")
	api.Use(authMiddleware)
	api.GET("/user", func(c echo.Context) error {
		c.Logger().Info("User endpoint hit")
		user := c.Get("user").(*models.User)
		return c.JSON(200, user)
	})
	r.Logger.Fatal(r.Start(cfg.ServerAddress))
}
