package email

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeEmail   NotificationType = "email"
	NotificationTypePush    NotificationType = "push"
	NotificationTypeSMS     NotificationType = "sms"
	NotificationTypeInApp   NotificationType = "in_app"
)

type NotificationStatus string

const (
	NotificationStatusPending   NotificationStatus = "pending"
	NotificationStatusSent      NotificationStatus = "sent"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusCancelled NotificationStatus = "cancelled"
)

type Notification struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Type      NotificationType   `json:"type"`
	Status    NotificationStatus `json:"status"`
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	ReadAt    *time.Time         `json:"read_at,omitempty"`
	SentAt    *time.Time         `json:"sent_at,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type NotificationService struct {
	emailService *EmailService
	channels     map[NotificationType]NotificationChannel
}

type NotificationChannel interface {
	Send(ctx context.Context, notification Notification) error
	GetType() NotificationType
}

func NewNotificationService(emailService *EmailService) *NotificationService {
	return &NotificationService{
		emailService: emailService,
		channels:     make(map[NotificationType]NotificationChannel),
	}
}

func (ns *NotificationService) RegisterChannel(channel NotificationChannel) {
	ns.channels[channel.GetType()] = channel
}

func (ns *NotificationService) Send(ctx context.Context, notification Notification) error {
	channel, exists := ns.channels[notification.Type]
	if !exists {
		return ErrChannelNotFound
	}

	return channel.Send(ctx, notification)
}

func (ns *NotificationService) SendMultiple(ctx context.Context, notifications []Notification) error {
	for _, notification := range notifications {
		if err := ns.Send(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}

type EmailNotificationChannel struct {
	emailService *EmailService
	fromEmail    string
}

func NewEmailNotificationChannel(emailService *EmailService, fromEmail string) *EmailNotificationChannel {
	return &EmailNotificationChannel{
		emailService: emailService,
		fromEmail:    fromEmail,
	}
}

func (enc *EmailNotificationChannel) Send(ctx context.Context, notification Notification) error {
	userEmail, ok := notification.Data["email"].(string)
	if !ok || userEmail == "" {
		return ErrInvalidRecipient
	}

	message := EmailMessage{
		To:      []string{userEmail},
		Subject: notification.Title,
		Body:    notification.Message,
	}

	if htmlBody, ok := notification.Data["html_body"].(string); ok {
		message.HTMLBody = htmlBody
	}

	return enc.emailService.SendAsync(message)
}

func (enc *EmailNotificationChannel) GetType() NotificationType {
	return NotificationTypeEmail
}

type InAppNotificationChannel struct {
	notifications []Notification
}

func NewInAppNotificationChannel() *InAppNotificationChannel {
	return &InAppNotificationChannel{
		notifications: make([]Notification, 0),
	}
}

func (ianc *InAppNotificationChannel) Send(ctx context.Context, notification Notification) error {
	ianc.notifications = append(ianc.notifications, notification)
	return nil
}

func (ianc *InAppNotificationChannel) GetType() NotificationType {
	return NotificationTypeInApp
}

func (ianc *InAppNotificationChannel) GetNotifications(userID uuid.UUID) []Notification {
	result := make([]Notification, 0)
	for _, notif := range ianc.notifications {
		if notif.UserID == userID {
			result = append(result, notif)
		}
	}
	return result
}

func (ianc *InAppNotificationChannel) MarkAsRead(notificationID uuid.UUID) error {
	for i, notif := range ianc.notifications {
		if notif.ID == notificationID {
			now := time.Now()
			ianc.notifications[i].ReadAt = &now
			return nil
		}
	}
	return ErrNotificationNotFound
}

var (
	ErrChannelNotFound      = fmt.Errorf("notification channel not found")
	ErrInvalidRecipient     = fmt.Errorf("invalid recipient")
	ErrNotificationNotFound = fmt.Errorf("notification not found")
)

func (ns *NotificationService) NotifyUser(ctx context.Context, userID uuid.UUID, notifType NotificationType, title, message string, data map[string]interface{}) error {
	notification := Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Status:    NotificationStatusPending,
		Title:     title,
		Message:   message,
		Data:      data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ns.Send(ctx, notification); err != nil {
		return err
	}

	now := time.Now()
	notification.SentAt = &now
	notification.Status = NotificationStatusSent

	return nil
}
