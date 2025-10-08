package services

import (
	"context"
	"errors"
	"time"

	"base-golang-restful/app/auth"
	appErrors "base-golang-restful/app/errors"
	"base-golang-restful/app/models"
	"base-golang-restful/app/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenRevoked  = errors.New("refresh token has been revoked")
	ErrTokenExpired  = errors.New("refresh token has expired")
	ErrMaxTokens     = errors.New("maximum number of active tokens reached")
)

type TokenService struct {
	*BaseService
	jwtManager     *auth.JWTManager
	tokenRepo      repository.RefreshTokenRepository
	maxTokensPerUser int
}

type TokenServiceConfig struct {
	DB               *gorm.DB
	JWTManager       *auth.JWTManager
	TokenRepo        repository.RefreshTokenRepository
	MaxTokensPerUser int
}

func NewTokenService(config TokenServiceConfig) *TokenService {
	if config.MaxTokensPerUser == 0 {
		config.MaxTokensPerUser = 5
	}

	return &TokenService{
		BaseService:      NewBaseService(config.DB),
		jwtManager:       config.JWTManager,
		tokenRepo:        config.TokenRepo,
		maxTokensPerUser: config.MaxTokensPerUser,
	}
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type RefreshTokenMetadata struct {
	IPAddress string
	UserAgent string
}

func (s *TokenService) GenerateTokenPair(ctx context.Context, userID uuid.UUID, email string, metadata *RefreshTokenMetadata) (*TokenPair, error) {
	count, err := s.tokenRepo.CountActiveTokensByUser(ctx, userID)
	if err != nil {
		return nil, appErrors.DatabaseError("failed to count active tokens")
	}

	if count >= int64(s.maxTokensPerUser) {
		return nil, appErrors.BadRequest("maximum number of active tokens reached")
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(userID, email)
	if err != nil {
		return nil, appErrors.InternalServer("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID, email)
	if err != nil {
		return nil, appErrors.InternalServer("failed to generate refresh token")
	}

	refreshTokenDuration := s.jwtManager.GetTokenDuration(auth.RefreshToken)
	expiresAt := time.Now().Add(refreshTokenDuration)

	tokenRecord := &models.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}

	if metadata != nil {
		tokenRecord.IPAddress = metadata.IPAddress
		tokenRecord.UserAgent = metadata.UserAgent
	}

	if err := s.tokenRepo.Create(ctx, tokenRecord); err != nil {
		return nil, appErrors.DatabaseError("failed to save refresh token")
	}

	accessTokenDuration := s.jwtManager.GetTokenDuration(auth.AccessToken)

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenDuration.Seconds()),
		TokenType:    "Bearer",
		ExpiresAt:    time.Now().Add(accessTokenDuration),
	}, nil
}

func (s *TokenService) RefreshToken(ctx context.Context, refreshToken string, metadata *RefreshTokenMetadata) (*TokenPair, error) {
	tokenRecord, err := s.tokenRepo.FindValidByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.Unauthorized("invalid refresh token")
		}
		return nil, appErrors.DatabaseError("failed to find refresh token")
	}

	if tokenRecord.IsRevoked {
		return nil, appErrors.Unauthorized("refresh token has been revoked")
	}

	if tokenRecord.IsExpired() {
		return nil, appErrors.Unauthorized("refresh token has expired")
	}

	if err := s.tokenRepo.RevokeToken(ctx, refreshToken); err != nil {
		return nil, appErrors.DatabaseError("failed to revoke old token")
	}

	return s.GenerateTokenPair(ctx, tokenRecord.UserID, tokenRecord.User.Email, metadata)
}

func (s *TokenService) RevokeToken(ctx context.Context, token string) error {
	tokenRecord, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appErrors.NotFound("refresh token not found")
		}
		return appErrors.DatabaseError("failed to find refresh token")
	}

	if tokenRecord.IsRevoked {
		return nil
	}

	if err := s.tokenRepo.RevokeToken(ctx, token); err != nil {
		return appErrors.DatabaseError("failed to revoke token")
	}

	return nil
}

func (s *TokenService) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		return appErrors.DatabaseError("failed to revoke user tokens")
	}
	return nil
}

func (s *TokenService) GetUserTokens(ctx context.Context, userID uuid.UUID) ([]models.RefreshToken, error) {
	tokens, err := s.tokenRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, appErrors.DatabaseError("failed to get user tokens")
	}
	return tokens, nil
}

func (s *TokenService) CleanupExpiredTokens(ctx context.Context) error {
	if err := s.tokenRepo.DeleteExpiredTokens(ctx); err != nil {
		return appErrors.DatabaseError("failed to cleanup expired tokens")
	}
	return nil
}

func (s *TokenService) ValidateRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	tokenRecord, err := s.tokenRepo.FindValidByToken(ctx, token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.Unauthorized("invalid refresh token")
		}
		return nil, appErrors.DatabaseError("failed to validate refresh token")
	}

	return tokenRecord, nil
}
