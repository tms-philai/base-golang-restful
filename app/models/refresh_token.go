package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Token     string         `gorm:"type:varchar(500);uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expires_at"`
	IsRevoked bool           `gorm:"default:false;not null;index" json:"is_revoked"`
	RevokedAt *time.Time     `json:"revoked_at,omitempty"`
	IPAddress string         `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent string         `gorm:"type:varchar(500)" json:"user_agent,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

func (r *RefreshToken) IsValid() bool {
	return !r.IsRevoked && !r.IsExpired()
}

func (r *RefreshToken) Revoke() {
	r.IsRevoked = true
	now := time.Now()
	r.RevokedAt = &now
}
