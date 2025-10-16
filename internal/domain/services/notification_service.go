package services

import (
	"base-gin/internal/domain/interfaces"
	"base-gin/internal/domain/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	emailService *EmailService
	observers    map[models.NotificationType]interfaces.Observer
}

// Implement Subject interface
func (ns *NotificationService) Attach(observer interfaces.Observer) {
	ns.observers[observer.GetType()] = observer
}

func (ns *NotificationService) Detach(observerType models.NotificationType) {
	delete(ns.observers, observerType)
}

func (ns *NotificationService) Notify(ctx context.Context, notification models.Notification) error {
	observer, exists := ns.observers[notification.Type]
	if !exists {
		return ErrChannelNotFound
	}
	return observer.Update(ctx, notification)
}

func NewNotificationService(emailService *EmailService) *NotificationService {
	return &NotificationService{
		emailService: emailService,
		observers:    make(map[models.NotificationType]interfaces.Observer),
	}
}

// Send sends a notification using Observer Pattern
func (ns *NotificationService) Send(ctx context.Context, notification models.Notification) error {
	return ns.Notify(ctx, notification)
}

// SendMultiple sends multiple notifications
func (ns *NotificationService) SendMultiple(ctx context.Context, notifications []models.Notification) error {
	for _, notification := range notifications {
		if err := ns.Send(ctx, notification); err != nil {
			return err
		}
	}
	return nil
}

// NotifyUser sends a notification to a user
func (ns *NotificationService) NotifyUser(ctx context.Context, userID uuid.UUID, notifType models.NotificationType, title, message string, data map[string]interface{}) error {
	notification := models.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Status:    models.NotificationStatusPending,
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
	notification.Status = models.NotificationStatusSent

	return nil
}

var (
	ErrChannelNotFound      = fmt.Errorf("notification channel not found")
	ErrInvalidRecipient     = fmt.Errorf("invalid recipient")
	ErrNotificationNotFound = fmt.Errorf("notification not found")
)
