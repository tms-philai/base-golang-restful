package services

import (
	"testing"
	"time"

	"base-golang-restful-app/config"
	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/stretchr/testify/assert"
)

// TestJWTServiceCreation tests JWT service creation
func TestJWTServiceCreation(t *testing.T) {
	// Setup
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key-for-jwt-testing",
		AccessTokenDuration:  3600,
		RefreshTokenDuration: 86400,
		Issuer:               "test-issuer",
	}
	jwtService := services.NewJWTService(cfg)

	// Test
	assert.NotNil(t, jwtService)
}

// TestJWTServiceTokenGeneration tests basic token generation
func TestJWTServiceTokenGeneration(t *testing.T) {
	// Setup
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key-for-jwt-testing-very-long-key",
		AccessTokenDuration:  3600,
		RefreshTokenDuration: 86400,
		Issuer:               "test-issuer",
	}
	jwtService := services.NewJWTService(cfg)

	testUser := &models.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      "user",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test token generation
	tokenPair, err := jwtService.GenerateTokenPair(testUser)
	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
}

// TestJWTServiceExtractTokenBasic tests basic token extraction
func TestJWTServiceExtractTokenBasic(t *testing.T) {
	// Setup
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key",
		AccessTokenDuration:  3600,
		RefreshTokenDuration: 86400,
		Issuer:               "test-issuer",
	}
	jwtService := services.NewJWTService(cfg)

	// Test successful extraction
	token := "sample.jwt.token"
	authHeader := "Bearer " + token

	extractedToken, err := jwtService.ExtractTokenFromHeader(authHeader)
	assert.NoError(t, err)
	assert.Equal(t, token, extractedToken)

	// Test empty header
	_, err = jwtService.ExtractTokenFromHeader("")
	assert.Error(t, err)
}
