package services

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt"
)

type TokenService struct {
	secretKey           string
	tokenValidityPeriod time.Duration
}

func NewTokenService(secretKey string, tokenValidity time.Duration) *TokenService {
	return &TokenService{secretKey: secretKey, tokenValidityPeriod: tokenValidity}
}

func (s *TokenService) GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.tokenValidityPeriod).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *TokenService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}
