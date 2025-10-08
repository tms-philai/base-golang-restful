package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTManager(t *testing.T) {
	tests := []struct {
		name   string
		config JWTConfig
	}{
		{
			name: "with custom config",
			config: JWTConfig{
				SecretKey:            "test-secret",
				AccessTokenDuration:  30 * time.Minute,
				RefreshTokenDuration: 14 * 24 * time.Hour,
				Issuer:               "test-issuer",
			},
		},
		{
			name: "with default values",
			config: JWTConfig{
				SecretKey: "test-secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewJWTManager(tt.config)

			assert.NotNil(t, manager)
			assert.Equal(t, tt.config.SecretKey, manager.secretKey)

			if tt.config.AccessTokenDuration == 0 {
				assert.Equal(t, 15*time.Minute, manager.accessTokenDuration)
			} else {
				assert.Equal(t, tt.config.AccessTokenDuration, manager.accessTokenDuration)
			}

			if tt.config.RefreshTokenDuration == 0 {
				assert.Equal(t, 7*24*time.Hour, manager.refreshTokenDuration)
			} else {
				assert.Equal(t, tt.config.RefreshTokenDuration, manager.refreshTokenDuration)
			}

			if tt.config.Issuer == "" {
				assert.Equal(t, "base-golang-restful", manager.issuer)
			} else {
				assert.Equal(t, tt.config.Issuer, manager.issuer)
			}
		})
	}
}

func TestJWTManager_GenerateAccessToken(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey: "test-secret",
	})

	userID := uuid.New()
	email := "test@example.com"

	token, err := manager.GenerateAccessToken(userID, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ValidateToken(token, AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, AccessToken, claims.TokenType)
}

func TestJWTManager_GenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey: "test-secret",
	})

	userID := uuid.New()
	email := "test@example.com"

	token, err := manager.GenerateRefreshToken(userID, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ValidateToken(token, RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, RefreshToken, claims.TokenType)
}

func TestJWTManager_ValidateToken(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey:           "test-secret",
		AccessTokenDuration: 1 * time.Second,
	})

	userID := uuid.New()
	email := "test@example.com"

	tests := []struct {
		name          string
		setupToken    func() string
		expectedType  TokenType
		expectError   bool
		expectedError error
	}{
		{
			name: "valid access token",
			setupToken: func() string {
				token, _ := manager.GenerateAccessToken(userID, email)
				return token
			},
			expectedType: AccessToken,
			expectError:  false,
		},
		{
			name: "valid refresh token",
			setupToken: func() string {
				token, _ := manager.GenerateRefreshToken(userID, email)
				return token
			},
			expectedType: RefreshToken,
			expectError:  false,
		},
		{
			name: "wrong token type",
			setupToken: func() string {
				token, _ := manager.GenerateAccessToken(userID, email)
				return token
			},
			expectedType:  RefreshToken,
			expectError:   true,
			expectedError: ErrInvalidToken,
		},
		{
			name: "invalid token",
			setupToken: func() string {
				return "invalid.token.string"
			},
			expectedType:  AccessToken,
			expectError:   true,
			expectedError: ErrInvalidToken,
		},
		{
			name: "expired token",
			setupToken: func() string {
				token, _ := manager.GenerateAccessToken(userID, email)
				time.Sleep(2 * time.Second)
				return token
			},
			expectedType:  AccessToken,
			expectError:   true,
			expectedError: ErrExpiredToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			claims, err := manager.ValidateToken(token, tt.expectedType)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, userID, claims.UserID)
				assert.Equal(t, email, claims.Email)
			}
		})
	}
}

func TestJWTManager_RefreshAccessToken(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey: "test-secret",
	})

	userID := uuid.New()
	email := "test@example.com"

	refreshToken, err := manager.GenerateRefreshToken(userID, email)
	assert.NoError(t, err)

	newAccessToken, err := manager.RefreshAccessToken(refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)

	claims, err := manager.ValidateToken(newAccessToken, AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
}

func TestJWTManager_RefreshAccessToken_InvalidToken(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey: "test-secret",
	})

	userID := uuid.New()
	email := "test@example.com"

	accessToken, err := manager.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	_, err = manager.RefreshAccessToken(accessToken)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTManager_GetTokenDuration(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey:            "test-secret",
		AccessTokenDuration:  30 * time.Minute,
		RefreshTokenDuration: 14 * 24 * time.Hour,
	})

	accessDuration := manager.GetTokenDuration(AccessToken)
	assert.Equal(t, 30*time.Minute, accessDuration)

	refreshDuration := manager.GetTokenDuration(RefreshToken)
	assert.Equal(t, 14*24*time.Hour, refreshDuration)
}

func TestJWTManager_ExtractClaims(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey: "test-secret",
	})

	userID := uuid.New()
	email := "test@example.com"

	token, err := manager.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	claims, err := manager.ExtractClaims(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, AccessToken, claims.TokenType)
}

func TestJWTManager_InvalidSignature(t *testing.T) {
	manager1 := NewJWTManager(JWTConfig{
		SecretKey: "secret-1",
	})

	manager2 := NewJWTManager(JWTConfig{
		SecretKey: "secret-2",
	})

	userID := uuid.New()
	email := "test@example.com"

	token, err := manager1.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	_, err = manager2.ValidateToken(token, AccessToken)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestJWTClaims_MissingFields(t *testing.T) {
	manager := NewJWTManager(JWTConfig{
		SecretKey: "test-secret",
	})

	claims := JWTClaims{
		UserID:    uuid.Nil,
		Email:     "",
		TokenType: AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(manager.secretKey))
	assert.NoError(t, err)

	_, err = manager.ValidateToken(tokenString, AccessToken)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingClaims)
}
