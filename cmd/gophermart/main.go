package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	mid "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
	r.Logger.SetLevel(log.DEBUG)
	r.HideBanner = true

	r.Use(mid.Recover())
	r.Use(middleware.Gzip())
	r.Use(middleware.Decompress())
	r.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Logger().Infof("Request Body: %v", string(reqBody))
		c.Logger().Infof("Response Body: %v", string(resBody))
	}))
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
