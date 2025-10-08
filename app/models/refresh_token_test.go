package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{
			name:      "not expired",
			expiresAt: time.Now().Add(1 * time.Hour),
			expected:  false,
		},
		{
			name:      "expired",
			expiresAt: time.Now().Add(-1 * time.Hour),
			expected:  true,
		},
		{
			name:      "just expired",
			expiresAt: time.Now().Add(-1 * time.Second),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := &RefreshToken{
				ExpiresAt: tt.expiresAt,
			}

			result := token.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRefreshToken_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		token     *RefreshToken
		expected  bool
	}{
		{
			name: "valid token",
			token: &RefreshToken{
				IsRevoked: false,
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			expected: true,
		},
		{
			name: "revoked token",
			token: &RefreshToken{
				IsRevoked: true,
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "expired token",
			token: &RefreshToken{
				IsRevoked: false,
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "revoked and expired",
			token: &RefreshToken{
				IsRevoked: true,
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.token.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRefreshToken_Revoke(t *testing.T) {
	token := &RefreshToken{
		ID:        uuid.New(),
		IsRevoked: false,
		RevokedAt: nil,
	}

	assert.False(t, token.IsRevoked)
	assert.Nil(t, token.RevokedAt)

	token.Revoke()

	assert.True(t, token.IsRevoked)
	assert.NotNil(t, token.RevokedAt)
	assert.WithinDuration(t, time.Now(), *token.RevokedAt, 1*time.Second)
}

func TestRefreshToken_BeforeCreate(t *testing.T) {
	token := &RefreshToken{
		UserID: uuid.New(),
		Token:  "test-token",
	}

	assert.Equal(t, uuid.Nil, token.ID)

	err := token.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, token.ID)
}

func TestRefreshToken_TableName(t *testing.T) {
	token := RefreshToken{}
	assert.Equal(t, "refresh_tokens", token.TableName())
}
