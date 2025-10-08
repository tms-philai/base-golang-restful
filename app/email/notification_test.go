package email

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotificationService_RegisterChannel(t *testing.T) {
	ns := NewNotificationService(nil)

	channel := NewInAppNotificationChannel()
	ns.RegisterChannel(channel)

	assert.Len(t, ns.channels, 1)
	assert.NotNil(t, ns.channels[NotificationTypeInApp])
}

func TestNotificationService_Send(t *testing.T) {
	ns := NewNotificationService(nil)

	channel := NewInAppNotificationChannel()
	ns.RegisterChannel(channel)

	notification := Notification{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Type:    NotificationTypeInApp,
		Status:  NotificationStatusPending,
		Title:   "Test Notification",
		Message: "This is a test",
	}

	err := ns.Send(context.Background(), notification)
	assert.NoError(t, err)
}

func TestNotificationService_SendMultiple(t *testing.T) {
	ns := NewNotificationService(nil)

	channel := NewInAppNotificationChannel()
	ns.RegisterChannel(channel)

	notifications := []Notification{
		{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Type:    NotificationTypeInApp,
			Title:   "Notification 1",
			Message: "Message 1",
		},
		{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Type:    NotificationTypeInApp,
			Title:   "Notification 2",
			Message: "Message 2",
		},
	}

	err := ns.SendMultiple(context.Background(), notifications)
	assert.NoError(t, err)
}

func TestInAppNotificationChannel_Send(t *testing.T) {
	channel := NewInAppNotificationChannel()

	notification := Notification{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Type:    NotificationTypeInApp,
		Title:   "Test",
		Message: "Test message",
	}

	err := channel.Send(context.Background(), notification)
	assert.NoError(t, err)

	notifications := channel.GetNotifications(notification.UserID)
	assert.Len(t, notifications, 1)
	assert.Equal(t, notification.Title, notifications[0].Title)
}

func TestInAppNotificationChannel_MarkAsRead(t *testing.T) {
	channel := NewInAppNotificationChannel()

	notificationID := uuid.New()
	notification := Notification{
		ID:      notificationID,
		UserID:  uuid.New(),
		Type:    NotificationTypeInApp,
		Title:   "Test",
		Message: "Test message",
	}

	err := channel.Send(context.Background(), notification)
	assert.NoError(t, err)

	err = channel.MarkAsRead(notificationID)
	assert.NoError(t, err)

	notifications := channel.GetNotifications(notification.UserID)
	assert.NotNil(t, notifications[0].ReadAt)
}

func TestInAppNotificationChannel_GetNotifications(t *testing.T) {
	channel := NewInAppNotificationChannel()

	userID1 := uuid.New()
	userID2 := uuid.New()

	notification1 := Notification{
		ID:      uuid.New(),
		UserID:  userID1,
		Type:    NotificationTypeInApp,
		Title:   "User 1 Notification",
		Message: "Message for user 1",
	}

	notification2 := Notification{
		ID:      uuid.New(),
		UserID:  userID2,
		Type:    NotificationTypeInApp,
		Title:   "User 2 Notification",
		Message: "Message for user 2",
	}

	channel.Send(context.Background(), notification1)
	channel.Send(context.Background(), notification2)

	user1Notifications := channel.GetNotifications(userID1)
	assert.Len(t, user1Notifications, 1)
	assert.Equal(t, "User 1 Notification", user1Notifications[0].Title)

	user2Notifications := channel.GetNotifications(userID2)
	assert.Len(t, user2Notifications, 1)
	assert.Equal(t, "User 2 Notification", user2Notifications[0].Title)
}

func TestNotificationService_NotifyUser(t *testing.T) {
	ns := NewNotificationService(nil)

	channel := NewInAppNotificationChannel()
	ns.RegisterChannel(channel)

	userID := uuid.New()
	data := map[string]interface{}{
		"key": "value",
	}

	err := ns.NotifyUser(
		context.Background(),
		userID,
		NotificationTypeInApp,
		"Test Title",
		"Test Message",
		data,
	)

	assert.NoError(t, err)

	notifications := channel.GetNotifications(userID)
	assert.Len(t, notifications, 1)
	assert.Equal(t, "Test Title", notifications[0].Title)
	assert.Equal(t, "Test Message", notifications[0].Message)
}

func TestNotification_Status(t *testing.T) {
	notification := Notification{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Type:   NotificationTypeEmail,
		Status: NotificationStatusPending,
	}

	assert.Equal(t, NotificationStatusPending, notification.Status)

	notification.Status = NotificationStatusSent
	now := time.Now()
	notification.SentAt = &now

	assert.Equal(t, NotificationStatusSent, notification.Status)
	assert.NotNil(t, notification.SentAt)
}
