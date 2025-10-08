package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmailClient(t *testing.T) {
	config := SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
		UseTLS:   true,
	}

	client := NewEmailClient(config)

	assert.NotNil(t, client)
	assert.NotNil(t, client.dialer)
	assert.Equal(t, "noreply@example.com", client.from)
}

func TestEmailClient_Send(t *testing.T) {
	t.Skip("Requires SMTP server")

	config := SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
		UseTLS:   true,
	}

	client := NewEmailClient(config)

	message := EmailMessage{
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Body:    "This is a test email",
	}

	err := client.Send(message)
	assert.NoError(t, err)
}

func TestEmailClient_SendSimple(t *testing.T) {
	t.Skip("Requires SMTP server")

	config := SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
	}

	client := NewEmailClient(config)

	err := client.SendSimple(
		[]string{"recipient@example.com"},
		"Test Subject",
		"Test Body",
	)

	assert.NoError(t, err)
}

func TestEmailClient_SendHTML(t *testing.T) {
	t.Skip("Requires SMTP server")

	config := SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "password",
		From:     "noreply@example.com",
	}

	client := NewEmailClient(config)

	htmlBody := "<html><body><h1>Test Email</h1><p>This is a test</p></body></html>"

	err := client.SendHTML(
		[]string{"recipient@example.com"},
		"Test HTML Email",
		htmlBody,
	)

	assert.NoError(t, err)
}

func TestEmailMessage_WithAttachment(t *testing.T) {
	message := EmailMessage{
		To:      []string{"recipient@example.com"},
		Subject: "Test with Attachment",
		Body:    "Please find the attachment",
		Attachments: []Attachment{
			{
				Filename: "test.txt",
				Content:  []byte("test content"),
				MimeType: "text/plain",
			},
		},
	}

	assert.Len(t, message.Attachments, 1)
	assert.Equal(t, "test.txt", message.Attachments[0].Filename)
}

func TestEmailMessage_WithCcAndBcc(t *testing.T) {
	message := EmailMessage{
		To:      []string{"recipient@example.com"},
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"bcc@example.com"},
		Subject: "Test Email",
		Body:    "Test Body",
	}

	assert.Len(t, message.To, 1)
	assert.Len(t, message.Cc, 1)
	assert.Len(t, message.Bcc, 1)
}

func TestEmailMessage_WithReplyTo(t *testing.T) {
	message := EmailMessage{
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Body:    "Test Body",
		ReplyTo: "reply@example.com",
	}

	assert.Equal(t, "reply@example.com", message.ReplyTo)
}

func TestEmailMessage_WithHeaders(t *testing.T) {
	message := EmailMessage{
		To:      []string{"recipient@example.com"},
		Subject: "Test Email",
		Body:    "Test Body",
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
			"X-Priority":      "1",
		},
	}

	assert.Len(t, message.Headers, 2)
	assert.Equal(t, "custom-value", message.Headers["X-Custom-Header"])
}
