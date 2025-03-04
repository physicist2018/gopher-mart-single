package main

import (
	"github.com/gin-gonic/gin"
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

	r := gin.Default()
	userRepo := repository.NewUserRepository(db)
	authService := authservice.NewAuthService("superkey", userRepo)
	authMiddleware := middlewares.JWTAuthMiddleware(authService)

	handlers := handlers.NewHandler(userRepo, authService)

	r.POST("/api/user/register", handlers.RegisterUser)
	r.POST("/api/user/login", handlers.LoginUser)
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
	}
	//r.GET("/api/user/:id", handlers.GetUserByID)
	r.Run(":8080")
}
