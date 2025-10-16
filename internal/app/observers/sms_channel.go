package observers

import (
	"base-gin/internal/domain/models"
	"context"
	"fmt"
)

type SMSChannel struct {
	apiKey    string
	apiSecret string
}

func NewSMSChannel(apiKey, apiSecret string) *SMSChannel {
	return &SMSChannel{
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

// Implement Observer interface
func (sc *SMSChannel) Update(ctx context.Context, notification models.Notification) error {
	phoneNumber, ok := notification.Data["phone"].(string)
	if !ok || phoneNumber == "" {
		return fmt.Errorf("phone number is required for SMS notification")
	}

	// In real implementation, you would call SMS API here
	fmt.Printf("SMS sent to %s: %s\n", phoneNumber, notification.Message)

	// Simulate API call
	// return sc.sendSMS(phoneNumber, notification.Message)
	return nil
}

func (sc *SMSChannel) GetType() models.NotificationType {
	return "sms" // You might want to add this to NotificationType enum
}

// Helper method to send actual SMS
func (sc *SMSChannel) sendSMS(phoneNumber, message string) error {
	// This is a placeholder - implement actual SMS API call
	// Example: Twilio, AWS SNS, etc.
	return nil
}
