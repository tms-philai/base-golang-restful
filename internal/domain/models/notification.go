package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeInApp NotificationType = "in_app"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

type Notification struct {
	ID        uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID              `json:"user_id" gorm:"type:uuid;not null;index"`
	Type      NotificationType       `json:"type" gorm:"type:varchar(20);not null;index"`
	Status    NotificationStatus     `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
	Title     string                 `json:"title" gorm:"type:varchar(255);not null"`
	Message   string                 `json:"message" gorm:"type:text;not null"`
	Data      map[string]interface{} `json:"data,omitempty" gorm:"type:jsonb"`
	ReadAt    *time.Time             `json:"read_at,omitempty" gorm:"type:timestamp"`
	SentAt    *time.Time             `json:"sent_at,omitempty" gorm:"type:timestamp"`
	CreatedAt time.Time              `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time              `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

type SendNotificationRequest struct {
	UserID  string                 `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type    NotificationType       `json:"type" binding:"required" example:"email"`
	Title   string                 `json:"title" binding:"required" example:"Welcome!"`
	Message string                 `json:"message" binding:"required" example:"Thank you for joining us"`
	Data    map[string]interface{} `json:"data" example:"{\"email\":\"user@example.com\",\"html_body\":\"<p>Thank you for joining us</p>\"}"`
}

type SendBulkNotificationRequest struct {
	Notifications []SendNotificationRequest `json:"notifications" binding:"required,min=1"`
}

type NotificationResponse struct {
	ID        uuid.UUID          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    uuid.UUID          `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type      NotificationType   `json:"type" example:"email"`
	Status    NotificationStatus `json:"status" example:"sent"`
	Title     string             `json:"title" example:"Welcome!"`
	Message   string             `json:"message" example:"Thank you for joining us"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
	ReadAt    *time.Time         `json:"read_at,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
}

type ListNotificationsResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int                    `json:"total" example:"10"`
	Unread        int                    `json:"unread" example:"5"`
}

type MarkAsReadRequest struct {
	NotificationID string `json:"notification_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type NotificationServiceStatusResponse struct {
	EmailChannel bool   `json:"email_channel" example:"true"`
	InAppChannel bool   `json:"in_app_channel" example:"true"`
	Status       string `json:"status" example:"operational"`
	Message      string `json:"message" example:"Notification service is running"`
}
