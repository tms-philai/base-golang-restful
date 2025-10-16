package repository

import (
	"base-gin/internal/domain/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	if err := r.db.WithContext(ctx).Create(notification).Error; err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&notification).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return &notification, nil
}

func (r *NotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get notifications by user: %w", err)
	}
	return notifications, nil
}

func (r *NotificationRepository) GetUnreadByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Notification, error) {
	var notifications []*models.Notification
	query := r.db.WithContext(ctx).Where("user_id = ? AND read_at IS NULL", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}
	return notifications, nil
}

func (r *NotificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	if err := r.db.WithContext(ctx).Save(notification).Error; err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ?", id).
		Update("read_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkAsSent(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"sent_at": gorm.Expr("NOW()"),
			"status":  models.NotificationStatusSent,
		}).Error; err != nil {
		return fmt.Errorf("failed to mark notification as sent: %w", err)
	}
	return nil
}

func (r *NotificationRepository) MarkAsFailed(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("id = ?", id).
		Update("status", models.NotificationStatusFailed).Error; err != nil {
		return fmt.Errorf("failed to mark notification as failed: %w", err)
	}
	return nil
}

func (r *NotificationRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count notifications: %w", err)
	}
	return int(count), nil
}

func (r *NotificationRepository) CountUnreadByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return int(count), nil
}

func (r *NotificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Notification{}).Error; err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetAllNotifications(ctx context.Context) ([]*models.Notification, error) {
	var notifications []*models.Notification
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("failed to get all notifications: %w", err)
	}
	return notifications, nil
}
