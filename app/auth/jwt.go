package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrMissingClaims    = errors.New("missing required claims")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type JWTClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
}

func NewJWTManager(config JWTConfig) *JWTManager {
	if config.AccessTokenDuration == 0 {
		config.AccessTokenDuration = 15 * time.Minute
	}
	if config.RefreshTokenDuration == 0 {
		config.RefreshTokenDuration = 7 * 24 * time.Hour
	}
	if config.Issuer == "" {
		config.Issuer = "base-golang-restful"
	}

	return &JWTManager{
		secretKey:            config.SecretKey,
		accessTokenDuration:  config.AccessTokenDuration,
		refreshTokenDuration: config.RefreshTokenDuration,
		issuer:               config.Issuer,
	}
}

func (m *JWTManager) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	return m.generateToken(userID, email, AccessToken, m.accessTokenDuration)
}

func (m *JWTManager) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	return m.generateToken(userID, email, RefreshToken, m.refreshTokenDuration)
}

func (m *JWTManager) generateToken(userID uuid.UUID, email string, tokenType TokenType, duration time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:    userID,
		Email:     email,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) ValidateToken(tokenString string, expectedType TokenType) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}

	if claims.UserID == uuid.Nil || claims.Email == "" {
		return nil, ErrMissingClaims
	}

	return claims, nil
}

func (m *JWTManager) RefreshAccessToken(refreshTokenString string) (string, error) {
	claims, err := m.ValidateToken(refreshTokenString, RefreshToken)
	if err != nil {
		return "", err
	}

	return m.GenerateAccessToken(claims.UserID, claims.Email)
}

func (m *JWTManager) GetTokenDuration(tokenType TokenType) time.Duration {
	if tokenType == AccessToken {
		return m.accessTokenDuration
	}
	return m.refreshTokenDuration
}

func (m *JWTManager) ExtractClaims(tokenString string) (*JWTClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
