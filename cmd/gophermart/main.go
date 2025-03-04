package main

import (
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/physicist2018/gopher-mart-single/internal/handlers"
	"github.com/physicist2018/gopher-mart-single/internal/middlewares"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/repository"
	"github.com/physicist2018/gopher-mart-single/internal/services/authservice"
	"gorm.io/gorm"
)

func main() {
	// Миграция модели User
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
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
	//r.GET("/api/user/:id", handlers.GetUserByID)
	r.Run(":8080")
}
