package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenService_GenerateToken(t *testing.T) {
	secretKey := "test-secret"
	tokenValidity := 1 * time.Hour
	service := NewTokenService(secretKey, tokenValidity)

	// Test case: Generate a valid token
	userID := 123
	token, err := service.GenerateToken(userID)
	assert.NoError(t, err, "Expected no error when generating token")
	assert.NotEmpty(t, token, "Generated token should not be empty")

	// Test case: Validate the generated token
	validatedUserID, err := service.ValidateToken(token)
	assert.NoError(t, err, "Expected no error when validating token")
	assert.Equal(t, userID, validatedUserID, "User ID should match the one used to generate the token")
}

func TestTokenService_ValidateToken_InvalidToken(t *testing.T) {
	secretKey := "test-secret"
	tokenValidity := 1 * time.Hour
	service := NewTokenService(secretKey, tokenValidity)

	// Test case: Validate an invalid token
	invalidToken := "invalid-token"
	_, err := service.ValidateToken(invalidToken)
	assert.Error(t, err, "Expected an error when validating an invalid token")
	//assert.Equal(t, "invalid token", err.Error(), "Error message should match")
}

func TestTokenService_ValidateToken_ExpiredToken(t *testing.T) {
	secretKey := "test-secret"
	tokenValidity := -1 * time.Hour // Set token validity to the past to simulate an expired token
	service := NewTokenService(secretKey, tokenValidity)

	// Test case: Generate an expired token
	userID := 123
	token, err := service.GenerateToken(userID)
	assert.NoError(t, err, "Expected no error when generating token")

	// Test case: Validate the expired token
	_, err = service.ValidateToken(token)
	assert.Error(t, err, "Expected an error when validating an expired token")
	assert.Contains(t, err.Error(), "Token is expired", "Error message should indicate token expiration")
}

func TestTokenService_ValidateToken_InvalidSecretKey(t *testing.T) {
	secretKey := "test-secret"
	tokenValidity := 1 * time.Hour
	service := NewTokenService(secretKey, tokenValidity)

	// Test case: Generate a valid token
	userID := 123
	token, err := service.GenerateToken(userID)
	assert.NoError(t, err, "Expected no error when generating token")

	// Test case: Validate the token with an incorrect secret key
	invalidService := NewTokenService("wrong-secret", tokenValidity)
	_, err = invalidService.ValidateToken(token)
	assert.Error(t, err, "Expected an error when validating token with incorrect secret key")
	assert.Contains(t, err.Error(), "signature is invalid", "Error message should indicate invalid signature")
}
