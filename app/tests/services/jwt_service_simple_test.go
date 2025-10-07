package services

import (
	"testing"
	"time"

	"base-golang-restful-app/config"
	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/stretchr/testify/assert"
)

// TestJWTServiceBasic tests basic JWT service functionality
func TestJWTServiceBasic(t *testing.T) {
	// Setup
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key-for-jwt-testing",
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

	// Test token validation
	claims, err := jwtService.ValidateToken(tokenPair.AccessToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, testUser.ID, claims.UserID)
	assert.Equal(t, testUser.Username, claims.Username)
	assert.Equal(t, testUser.Email, claims.Email)
	assert.Equal(t, testUser.Role, claims.Role)
}

// TestJWTServiceExtractToken tests token extraction from header
func TestJWTServiceExtractToken(t *testing.T) {
	// Setup
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key",
		AccessTokenDuration:  3600,
		RefreshTokenDuration: 86400,
		Issuer:               "test-issuer",
	}
	jwtService := services.NewJWTService(cfg)

	// Test successful extraction
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token"
	authHeader := "Bearer " + token

	extractedToken, err := jwtService.ExtractTokenFromHeader(authHeader)
	assert.NoError(t, err)
	assert.Equal(t, token, extractedToken)

	// Test empty header
	_, err = jwtService.ExtractTokenFromHeader("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "authorization header is required")

	// Test missing Bearer
	_, err = jwtService.ExtractTokenFromHeader("InvalidHeader")
	assert.Error(t, err)
}

// TestJWTServiceInvalidToken tests validation of invalid tokens
func TestJWTServiceInvalidToken(t *testing.T) {
	// Setup
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key",
		AccessTokenDuration:  3600,
		RefreshTokenDuration: 86400,
		Issuer:               "test-issuer",
	}
	jwtService := services.NewJWTService(cfg)

	// Test invalid token
	claims, err := jwtService.ValidateToken("invalid-token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}
