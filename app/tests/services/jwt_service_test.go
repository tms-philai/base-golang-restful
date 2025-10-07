package services

import (
	"testing"
	"time"

	"base-golang-restful-app/config"
	"base-golang-restful-app/models"
	"base-golang-restful-app/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// JWTServiceTestSuite defines the test suite for JWTService
type JWTServiceTestSuite struct {
	suite.Suite
	jwtService *services.JWTService
	testUser   *models.User
}

// SetupTest sets up the test environment before each test
func (suite *JWTServiceTestSuite) SetupTest() {
	cfg := &config.JWTConfig{
		SecretKey:            "test-secret-key-for-jwt-testing",
		AccessTokenDuration:  3600,    // 1 hour
		RefreshTokenDuration: 86400,   // 24 hours
		Issuer:               "test-issuer",
	}
	suite.jwtService = services.NewJWTService(cfg)

	// Create test user
	suite.testUser = &models.User{
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
}

// TestGenerateTokenPairSuccess tests successful token pair generation
func (suite *JWTServiceTestSuite) TestGenerateTokenPairSuccess() {
	// Act
	tokenPair, err := suite.jwtService.GenerateTokenPair(suite.testUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), tokenPair)
	assert.NotEmpty(suite.T(), tokenPair.AccessToken)
	assert.NotEmpty(suite.T(), tokenPair.RefreshToken)
	assert.Equal(suite.T(), int64(3600), tokenPair.ExpiresIn)
}

// TestValidateTokenSuccess tests successful token validation
func (suite *JWTServiceTestSuite) TestValidateTokenSuccess() {
	// Arrange
	tokenPair, _ := suite.jwtService.GenerateTokenPair(suite.testUser)

	// Act
	claims, err := suite.jwtService.ValidateToken(tokenPair.AccessToken)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), suite.testUser.ID, claims.UserID)
	assert.Equal(suite.T(), suite.testUser.Username, claims.Username)
	assert.Equal(suite.T(), suite.testUser.Email, claims.Email)
	assert.Equal(suite.T(), suite.testUser.Role, claims.Role)
	assert.Equal(suite.T(), "test-issuer", claims.Issuer)
}

// TestValidateTokenInvalid tests validation of invalid token
func (suite *JWTServiceTestSuite) TestValidateTokenInvalid() {
	// Act
	claims, err := suite.jwtService.ValidateToken("invalid-token")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "invalid token")
}

// TestValidateTokenExpired tests validation of expired token
func (suite *JWTServiceTestSuite) TestValidateTokenExpired() {
	// Arrange - Create JWT service with very short expiration
	shortCfg := &config.JWTConfig{
		SecretKey:            "test-secret-key",
		AccessTokenDuration:  1, // 1 second
		RefreshTokenDuration: 86400,
		Issuer:               "test-issuer",
	}
	shortJWTService := services.NewJWTService(shortCfg)

	tokenPair, _ := shortJWTService.GenerateTokenPair(suite.testUser)
	
	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Act
	claims, err := shortJWTService.ValidateToken(tokenPair.AccessToken)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "token is expired")
}

// TestRefreshTokenSuccess tests successful token refresh
func (suite *JWTServiceTestSuite) TestRefreshTokenSuccess() {
	// Arrange
	tokenPair, _ := suite.jwtService.GenerateTokenPair(suite.testUser)

	// Act
	newTokenPair, err := suite.jwtService.RefreshToken(tokenPair.RefreshToken, suite.testUser)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), newTokenPair)
	assert.NotEmpty(suite.T(), newTokenPair.AccessToken)
	assert.NotEmpty(suite.T(), newTokenPair.RefreshToken)
	assert.NotEqual(suite.T(), tokenPair.AccessToken, newTokenPair.AccessToken)
	assert.NotEqual(suite.T(), tokenPair.RefreshToken, newTokenPair.RefreshToken)
}

// TestRefreshTokenInvalid tests refresh with invalid refresh token
func (suite *JWTServiceTestSuite) TestRefreshTokenInvalid() {
	// Act
	newTokenPair, err := suite.jwtService.RefreshToken("invalid-refresh-token", suite.testUser)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), newTokenPair)
	assert.Contains(suite.T(), err.Error(), "invalid refresh token")
}

// TestRefreshTokenWrongType tests refresh with access token instead of refresh token
func (suite *JWTServiceTestSuite) TestRefreshTokenWrongType() {
	// Arrange
	tokenPair, _ := suite.jwtService.GenerateTokenPair(suite.testUser)

	// Act - Try to refresh using access token instead of refresh token
	newTokenPair, err := suite.jwtService.RefreshToken(tokenPair.AccessToken, suite.testUser)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), newTokenPair)
	assert.Contains(suite.T(), err.Error(), "invalid token type")
}

// TestExtractTokenFromHeaderSuccess tests successful token extraction from header
func (suite *JWTServiceTestSuite) TestExtractTokenFromHeaderSuccess() {
	// Arrange
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token"
	authHeader := "Bearer " + token

	// Act
	extractedToken, err := suite.jwtService.ExtractTokenFromHeader(authHeader)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), token, extractedToken)
}

// TestExtractTokenFromHeaderMissingBearer tests extraction without Bearer prefix
func (suite *JWTServiceTestSuite) TestExtractTokenFromHeaderMissingBearer() {
	// Arrange
	authHeader := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token"

	// Act
	extractedToken, err := suite.jwtService.ExtractTokenFromHeader(authHeader)

	// Assert
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), extractedToken)
	assert.Contains(suite.T(), err.Error(), "authorization header must start with Bearer")
}

// TestExtractTokenFromHeaderEmpty tests extraction from empty header
func (suite *JWTServiceTestSuite) TestExtractTokenFromHeaderEmpty() {
	// Act
	extractedToken, err := suite.jwtService.ExtractTokenFromHeader("")

	// Assert
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), extractedToken)
	assert.Contains(suite.T(), err.Error(), "authorization header is required")
}

// TestExtractTokenFromHeaderOnlyBearer tests extraction with only "Bearer" in header
func (suite *JWTServiceTestSuite) TestExtractTokenFromHeaderOnlyBearer() {
	// Act
	extractedToken, err := suite.jwtService.ExtractTokenFromHeader("Bearer")

	// Assert
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), extractedToken)
	assert.Contains(suite.T(), err.Error(), "authorization header format must be Bearer {token}")
}

// TestTokenClaimsConsistency tests that generated and validated claims are consistent
func (suite *JWTServiceTestSuite) TestTokenClaimsConsistency() {
	// Arrange
	adminUser := &models.User{
		ID:        "admin-456",
		Username:  "admin",
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Role:      "admin",
		IsActive:  true,
	}

	// Act
	tokenPair, err := suite.jwtService.GenerateTokenPair(adminUser)
	assert.NoError(suite.T(), err)

	claims, err := suite.jwtService.ValidateToken(tokenPair.AccessToken)
	assert.NoError(suite.T(), err)

	// Assert
	assert.Equal(suite.T(), adminUser.ID, claims.UserID)
	assert.Equal(suite.T(), adminUser.Username, claims.Username)
	assert.Equal(suite.T(), adminUser.Email, claims.Email)
	assert.Equal(suite.T(), adminUser.Role, claims.Role)
	assert.True(suite.T(), claims.ExpiresAt.After(time.Now()))
	assert.True(suite.T(), claims.IssuedAt.Before(time.Now().Add(time.Second)))
}

// TestRefreshTokenClaimsValidation tests that refresh token has correct claims
func (suite *JWTServiceTestSuite) TestRefreshTokenClaimsValidation() {
	// Arrange
	tokenPair, _ := suite.jwtService.GenerateTokenPair(suite.testUser)

	// Act - Validate refresh token claims
	refreshClaims, err := suite.jwtService.ValidateToken(tokenPair.RefreshToken)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), refreshClaims)
	assert.Equal(suite.T(), suite.testUser.ID, refreshClaims.UserID)
	assert.True(suite.T(), refreshClaims.ExpiresAt.After(time.Now().Add(time.Hour*23))) // Should expire in ~24 hours
}

// TestJWTServiceTestSuite runs the test suite
func TestJWTServiceTestSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))
}
