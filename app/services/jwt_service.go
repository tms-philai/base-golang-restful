package services

import (
	"errors"
	"time"

	"base-golang-restful-app/config"
	"base-golang-restful-app/models"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token operations
type JWTService struct {
	config *config.JWTConfig
}

// NewJWTService creates a new JWT service
func NewJWTService(cfg *config.JWTConfig) *JWTService {
	return &JWTService{
		config: cfg,
	}
}

// GenerateTokenPair generates access and refresh token pair for a user
func (s *JWTService) GenerateTokenPair(user *models.User) (*models.TokenPair, error) {
	// Generate access token
	accessToken, err := s.generateToken(user, "access", s.config.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateToken(user, "refresh", s.config.RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessTokenDuration.Seconds()),
	}, nil
}

// generateToken generates a JWT token for a user
func (s *JWTService) generateToken(user *models.User, tokenType string, duration time.Duration) (string, error) {
	claims := models.NewJWTClaims(user, tokenType, duration)
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new access token from a valid refresh token
func (s *JWTService) RefreshToken(refreshTokenString string, user *models.User) (*models.TokenPair, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Verify token belongs to the user
	if claims.UserID != user.ID {
		return nil, errors.New("token does not belong to user")
	}

	// Generate new token pair
	return s.GenerateTokenPair(user)
}

// ExtractTokenFromHeader extracts token from Authorization header
func (s *JWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Check if header starts with "Bearer "
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	token := authHeader[len(bearerPrefix):]
	if token == "" {
		return "", errors.New("token is required")
	}

	return token, nil
}
