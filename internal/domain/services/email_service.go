package services

import (
	"base-gin/internal/pkg/utils"
	"context"
	"fmt"
	"sync"
)

type EmailService struct {
	client    *utils.EmailClient
	queue     chan utils.EmailMessage
	workers   int
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	templates *utils.TemplateManager
}

type EmailServiceConfig struct {
	Client    *utils.EmailClient
	Workers   int
	QueueSize int
	Templates *utils.TemplateManager
}

func NewEmailService(config EmailServiceConfig) *EmailService {
	if config.Workers == 0 {
		config.Workers = 5
	}
	if config.QueueSize == 0 {
		config.QueueSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &EmailService{
		client:    config.Client,
		queue:     make(chan utils.EmailMessage, config.QueueSize),
		workers:   config.Workers,
		ctx:       ctx,
		cancel:    cancel,
		templates: config.Templates,
	}

	service.startWorkers()

	return service
}

func (s *EmailService) startWorkers() {
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *EmailService) worker(id int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case message, ok := <-s.queue:
			if !ok {
				return
			}

			if err := s.client.Send(message); err != nil {
				fmt.Printf("Worker %d: Failed to send email to %v: %v\n", id, message.To, err)
			}
		}
	}
}

func (s *EmailService) SendAsync(message utils.EmailMessage) error {
	select {
	case s.queue <- message:
		return nil
	default:
		return fmt.Errorf("email queue is full")
	}
}

func (s *EmailService) SendSync(message utils.EmailMessage) error {
	return s.client.Send(message)
}

func (s *EmailService) SendTemplate(templateName string, to []string, data interface{}) error {
	if s.templates == nil {
		return fmt.Errorf("template manager not configured")
	}

	rendered, err := s.templates.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return s.SendAsync(utils.EmailMessage{
		To:       to,
		Subject:  rendered.Subject,
		HTMLBody: rendered.HTMLBody,
		Body:     rendered.TextBody,
	})
}

func (s *EmailService) SendTemplateSync(templateName string, to []string, data interface{}) error {
	if s.templates == nil {
		return fmt.Errorf("template manager not configured")
	}

	rendered, err := s.templates.Render(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return s.SendSync(utils.EmailMessage{
		To:       to,
		Subject:  rendered.Subject,
		HTMLBody: rendered.HTMLBody,
		Body:     rendered.TextBody,
	})
}

func (s *EmailService) SendBulk(messages []utils.EmailMessage) error {
	for _, message := range messages {
		if err := s.SendAsync(message); err != nil {
			return err
		}
	}
	return nil
}

func (s *EmailService) QueueSize() int {
	return len(s.queue)
}

func (s *EmailService) Shutdown() {
	s.cancel()
	close(s.queue)
	s.wg.Wait()
}

func (s *EmailService) TestConnection() error {
	return s.client.TestConnection()
}
