package models

type SendEmailRequest struct {
	To      []string `json:"to" binding:"required" example:"user@example.com"`
	Subject string   `json:"subject" binding:"required" example:"Welcome to Our Platform"`
	Body    string   `json:"body" binding:"required" example:"Thank you for joining us!"`
	IsHTML  bool     `json:"is_html" example:"false"`
}

type SendEmailWithTemplateRequest struct {
	To           []string               `json:"to" binding:"required" example:"user@example.com"`
	TemplateName string                 `json:"template_name" binding:"required" example:"welcome"`
	Data         map[string]interface{} `json:"data" binding:"required"`
}

type SendBulkEmailRequest struct {
	Emails []SendEmailRequest `json:"emails" binding:"required,min=1"`
}

type EmailResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Email sent successfully"`
	QueueID string `json:"queue_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type EmailStatusResponse struct {
	QueueSize int    `json:"queue_size" example:"5"`
	Status    string `json:"status" example:"operational"`
	Message   string `json:"message" example:"Email service is running"`
}

type TestEmailConnectionResponse struct {
	Success   bool   `json:"success" example:"true"`
	Message   string `json:"message" example:"Email connection successful"`
	SMTPHost  string `json:"smtp_host" example:"smtp.gmail.com"`
	SMTPPort  int    `json:"smtp_port" example:"587"`
	Connected bool   `json:"connected" example:"true"`
}
