package email

import (
	"crypto/tls"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailClient struct {
	dialer *gomail.Dialer
	from   string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	UseTLS   bool
}

func NewEmailClient(config SMTPConfig) *EmailClient {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	
	if config.UseTLS {
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         config.Host,
		}
	}

	dialer.Timeout = 10 * time.Second

	return &EmailClient{
		dialer: dialer,
		from:   config.From,
	}
}

type EmailMessage struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	HTMLBody    string
	Attachments []Attachment
	ReplyTo     string
	Headers     map[string]string
}

type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

func (c *EmailClient) Send(message EmailMessage) error {
	m := gomail.NewMessage()

	m.SetHeader("From", c.from)
	m.SetHeader("To", message.To...)

	if len(message.Cc) > 0 {
		m.SetHeader("Cc", message.Cc...)
	}

	if len(message.Bcc) > 0 {
		m.SetHeader("Bcc", message.Bcc...)
	}

	m.SetHeader("Subject", message.Subject)

	if message.ReplyTo != "" {
		m.SetHeader("Reply-To", message.ReplyTo)
	}

	for key, value := range message.Headers {
		m.SetHeader(key, value)
	}

	if message.HTMLBody != "" {
		m.SetBody("text/html", message.HTMLBody)
		if message.Body != "" {
			m.AddAlternative("text/plain", message.Body)
		}
	} else {
		m.SetBody("text/plain", message.Body)
	}

	for _, attachment := range message.Attachments {
		m.Attach(attachment.Filename, gomail.SetCopyFunc(func(w gomail.Writer) error {
			_, err := w.Write(attachment.Content)
			return err
		}))
	}

	if err := c.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (c *EmailClient) SendSimple(to []string, subject, body string) error {
	return c.Send(EmailMessage{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

func (c *EmailClient) SendHTML(to []string, subject, htmlBody string) error {
	return c.Send(EmailMessage{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody,
	})
}

func (c *EmailClient) SendWithAttachment(to []string, subject, body string, attachments []Attachment) error {
	return c.Send(EmailMessage{
		To:          to,
		Subject:     subject,
		Body:        body,
		Attachments: attachments,
	})
}

func (c *EmailClient) TestConnection() error {
	d := c.dialer
	
	conn, err := d.Dial()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	return nil
}
