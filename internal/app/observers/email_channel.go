package observers

import (
	"base-gin/internal/domain/interfaces"
	"base-gin/internal/domain/models"
	"base-gin/internal/domain/services"
	"base-gin/internal/pkg/utils"
	"context"
	"fmt"
)

type EmailChannel struct {
	emailService *services.EmailService
	fromEmail    string
}

// Ensure EmailChannel implements interfaces.Observer
var _ interfaces.Observer = (*EmailChannel)(nil)

func NewEmailChannel(emailService *services.EmailService, fromEmail string) *EmailChannel {
	return &EmailChannel{
		emailService: emailService,
		fromEmail:    fromEmail,
	}
}

// Implement interfaces.Observer interface
func (ec *EmailChannel) Update(ctx context.Context, notification models.Notification) error {
	userEmail, ok := notification.Data["email"].(string)
	if !ok || userEmail == "" {
		return fmt.Errorf("invalid recipient")
	}

	message := utils.EmailMessage{
		To:      []string{userEmail},
		Subject: notification.Title,
		Body:    notification.Message,
	}

	if htmlBody, ok := notification.Data["html_body"].(string); ok {
		message.HTMLBody = htmlBody
	}

	return ec.emailService.SendAsync(message)
}

func (ec *EmailChannel) GetType() models.NotificationType {
	return models.NotificationTypeEmail
}
