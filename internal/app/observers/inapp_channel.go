package observers

import (
	"base-gin/internal/domain/interfaces"
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/repository"
	"context"

	"github.com/google/uuid"
)

type InAppChannel struct {
	notificationRepo *repository.NotificationRepository
}

// Ensure InAppChannel implements interfaces.Observer
var _ interfaces.Observer = (*InAppChannel)(nil)

func NewInAppChannel(notificationRepo *repository.NotificationRepository) *InAppChannel {
	return &InAppChannel{
		notificationRepo: notificationRepo,
	}
}

// Implement interfaces.Observer interface
func (iac *InAppChannel) Update(ctx context.Context, notification models.Notification) error {
	// Save notification to database
	return iac.notificationRepo.Create(ctx, &notification)
}

func (iac *InAppChannel) GetType() models.NotificationType {
	return models.NotificationTypeInApp
}

// Additional methods for in-app specific functionality
func (iac *InAppChannel) GetNotifications(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	return iac.notificationRepo.GetByUserID(ctx, userID, limit, offset)
}

func (iac *InAppChannel) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	return iac.notificationRepo.MarkAsRead(ctx, notificationID)
}
