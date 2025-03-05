package authservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/physicist2018/gopher-mart-single/internal/models"
	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
	"github.com/physicist2018/gopher-mart-single/internal/ports/repository"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	secretKey      string
	userRepository repository.UserRepository // Зависимость от UserRepository
}

func NewAuthService(secretKey string, userRepository repository.UserRepository) authservice.AuthService {
	return &authService{secretKey: secretKey, userRepository: userRepository}
}

func (s *authService) Register(ctx context.Context, login, password string) (*models.User, error) {
	// Регистрация пользователя
	existingUser, err := s.userRepository.GetUserByLogin(login)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if existingUser != nil {
		return nil, authservice.ErrUserAlreadyExists
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем нового пользователя
	user := &models.User{
		Login:    login,
		Password: string(hashedPassword),
	}

	// Сохраняем пользователя в базе данных
	if err := s.userRepository.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	log.Println("register Secret key: ", s.secretKey)
	return user, nil
}

func (s *authService) Login(ctx context.Context, login, password string) (string, error) {
	// Аутентификация пользователя
	// Получаем пользователя по логину
	user, err := s.userRepository.GetUserByLogin(login)
	if err != nil {

		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return "", authservice.ErrUserNotFound
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", authservice.ErrInvalidCredentials
	}

	log.Println("login Secret key: ", s.secretKey)
	// Генерируем JWT-токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Валидация токена
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Проверяем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Получаем ID пользователя из токена
	userID := claims["sub"].(float64)

	// Получаем пользователя по ID
	user, err := s.userRepository.GetUserByID(uint(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
