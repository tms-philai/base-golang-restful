package services

import (
	"context"
	"testing"
	"time"

	"base-golang-restful-app/auth"
	"base-golang-restful-app/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, entity *models.RefreshToken) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.RefreshToken, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindAll(ctx context.Context) ([]models.RefreshToken, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Update(ctx context.Context, entity *models.RefreshToken) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.RefreshToken, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) FindValidByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) CountActiveTokensByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func TestTokenService_GenerateTokenPair(t *testing.T) {
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		SecretKey: "test-secret",
	})

	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		JWTManager:       jwtManager,
		TokenRepo:        mockRepo,
		MaxTokensPerUser: 5,
	})

	userID := uuid.New()
	email := "test@example.com"
	metadata := &RefreshTokenMetadata{
		IPAddress: "127.0.0.1",
		UserAgent: "Test Agent",
	}

	mockRepo.On("CountActiveTokensByUser", mock.Anything, userID).Return(int64(0), nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	tokenPair, err := service.GenerateTokenPair(context.Background(), userID, email, metadata)

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	assert.Equal(t, "Bearer", tokenPair.TokenType)
	assert.Greater(t, tokenPair.ExpiresIn, int64(0))

	mockRepo.AssertExpectations(t)
}

func TestTokenService_GenerateTokenPair_MaxTokensReached(t *testing.T) {
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		SecretKey: "test-secret",
	})

	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		JWTManager:       jwtManager,
		TokenRepo:        mockRepo,
		MaxTokensPerUser: 5,
	})

	userID := uuid.New()
	email := "test@example.com"

	mockRepo.On("CountActiveTokensByUser", mock.Anything, userID).Return(int64(5), nil)

	tokenPair, err := service.GenerateTokenPair(context.Background(), userID, email, nil)

	assert.Error(t, err)
	assert.Nil(t, tokenPair)

	mockRepo.AssertExpectations(t)
}

func TestTokenService_RefreshToken(t *testing.T) {
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		SecretKey: "test-secret",
	})

	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		JWTManager:       jwtManager,
		TokenRepo:        mockRepo,
		MaxTokensPerUser: 5,
	})

	userID := uuid.New()
	email := "test@example.com"
	oldToken := "old-refresh-token"

	tokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     oldToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
		User: models.User{
			ID:    userID,
			Email: email,
		},
	}

	mockRepo.On("FindValidByToken", mock.Anything, oldToken).Return(tokenRecord, nil)
	mockRepo.On("RevokeToken", mock.Anything, oldToken).Return(nil)
	mockRepo.On("CountActiveTokensByUser", mock.Anything, userID).Return(int64(1), nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	tokenPair, err := service.RefreshToken(context.Background(), oldToken, nil)

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)

	mockRepo.AssertExpectations(t)
}

func TestTokenService_RevokeToken(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		TokenRepo: mockRepo,
	})

	token := "test-token"
	tokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		Token:     token,
		IsRevoked: false,
	}

	mockRepo.On("FindByToken", mock.Anything, token).Return(tokenRecord, nil)
	mockRepo.On("RevokeToken", mock.Anything, token).Return(nil)

	err := service.RevokeToken(context.Background(), token)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTokenService_RevokeAllUserTokens(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		TokenRepo: mockRepo,
	})

	userID := uuid.New()

	mockRepo.On("RevokeAllUserTokens", mock.Anything, userID).Return(nil)

	err := service.RevokeAllUserTokens(context.Background(), userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTokenService_GetUserTokens(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		TokenRepo: mockRepo,
	})

	userID := uuid.New()
	expectedTokens := []models.RefreshToken{
		{ID: uuid.New(), UserID: userID},
		{ID: uuid.New(), UserID: userID},
	}

	mockRepo.On("FindByUserID", mock.Anything, userID).Return(expectedTokens, nil)

	tokens, err := service.GetUserTokens(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, tokens, 2)
	mockRepo.AssertExpectations(t)
}

func TestTokenService_CleanupExpiredTokens(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		TokenRepo: mockRepo,
	})

	mockRepo.On("DeleteExpiredTokens", mock.Anything).Return(nil)

	err := service.CleanupExpiredTokens(context.Background())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTokenService_ValidateRefreshToken(t *testing.T) {
	mockRepo := new(MockRefreshTokenRepository)
	service := NewTokenService(TokenServiceConfig{
		TokenRepo: mockRepo,
	})

	token := "valid-token"
	tokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
	}

	mockRepo.On("FindValidByToken", mock.Anything, token).Return(tokenRecord, nil)

	result, err := service.ValidateRefreshToken(context.Background(), token)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, token, result.Token)
	mockRepo.AssertExpectations(t)
}
