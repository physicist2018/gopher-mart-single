package authservice

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	db "github.com/physicist2018/gopher-mart-single/internal/database/db/postgres"
	"github.com/physicist2018/gopher-mart-single/internal/ports/authservice"
	"github.com/physicist2018/gopher-mart-single/internal/ports/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService интерфейс для регистрации, аутентификации и валидации
type AuthService interface {
	// Register регистрирует нового пользователя.
	Register(ctx context.Context, login, password string) (*db.User, error)

	// Login выполняет аутентификацию пользователя и возвращает токен.
	Login(ctx context.Context, login, password string) (string, error)

	// ValidateToken проверяет токен и возвращает информацию о пользователе.
	ValidateToken(ctx context.Context, token string) (*db.User, error)
}

// authService реализация AuthService
type authService struct {
	secretKey string
	userRepo  repository.UserRepository // Ваш репозиторий для работы с пользователями
}

// Claims структура для JWT claims
type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

// NewAuthService конструктор для authService
func NewAuthService(secretKey string, userRepo repository.UserRepository) AuthService {
	return &authService{
		secretKey: secretKey,
		userRepo:  userRepo,
	}
}

// Register регистрирует нового пользователя
func (as *authService) Register(ctx context.Context, login, password string) (*db.User, error) {

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Save user's username and hashed password
	dbuser, err := as.userRepo.CreateUser(ctx, db.CreateUserParams{
		Login:    login,
		Password: string(hashedPassword),
	})

	return &dbuser, err
}

// Login выполняет аутентификацию пользователя
func (as *authService) Login(ctx context.Context, login, password string) (string, error) {
	// Поиск пользователя по логину
	user, err := as.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", authservice.ErrInvalidCredentials
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", authservice.ErrInvalidCredentials
	}

	// Генерация JWT токена
	token, err := as.generateToken(user.Login)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken проверяет JWT токен и возвращает информацию о пользователе
func (as *authService) ValidateToken(ctx context.Context, tokenString string) (*db.User, error) {
	claims := &Claims{}

	// Парсим и валидируем токен
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(as.secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, authservice.ErrInvalidCredentials
	}

	// Поиск пользователя по логину
	user, err := as.userRepo.GetUserByLogin(ctx, claims.Login)
	if err != nil {
		return nil, authservice.ErrUserNotFound
	}

	return &user, nil
}

// generateToken генерирует JWT токен
func (as *authService) generateToken(login string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Время жизни токена — 24 часа
	claims := &Claims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем его с помощью секретного ключа
	return token.SignedString([]byte(as.secretKey))
}
