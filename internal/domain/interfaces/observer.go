package interfaces

import (
	"base-gin/internal/domain/models"
	"context"
)

// Observer defines the interface for notification observers
type Observer interface {
	Update(ctx context.Context, notification models.Notification) error
	GetType() models.NotificationType
}
