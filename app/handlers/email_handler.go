package handlers

import (
	"net/http"

	"base-golang-restful-app/email"
	"base-golang-restful-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmailHandler struct {
	emailService *email.EmailService
}

func NewEmailHandler(emailService *email.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// SendEmail godoc
// @Summary Send a simple email
// @Description Send a simple text or HTML email to one or more recipients
// @Tags Email
// @Accept json
// @Produce json
// @Param request body models.SendEmailRequest true "Email details"
// @Success 200 {object} models.EmailResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/email/send [post]
// @Security BearerAuth
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req models.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	message := email.EmailMessage{
		To:      req.To,
		Subject: req.Subject,
	}

	if req.IsHTML {
		message.HTMLBody = req.Body
	} else {
		message.Body = req.Body
	}

	if err := h.emailService.SendAsync(message); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to send email",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.EmailResponse{
		Success: true,
		Message: "Email queued successfully",
		QueueID: uuid.New().String(),
	})
}

// SendEmailSync godoc
// @Summary Send an email synchronously
// @Description Send an email synchronously and wait for confirmation
// @Tags Email
// @Accept json
// @Produce json
// @Param request body models.SendEmailRequest true "Email details"
// @Success 200 {object} models.EmailResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/email/send-sync [post]
// @Security BearerAuth
func (h *EmailHandler) SendEmailSync(c *gin.Context) {
	var req models.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	message := email.EmailMessage{
		To:      req.To,
		Subject: req.Subject,
	}

	if req.IsHTML {
		message.HTMLBody = req.Body
	} else {
		message.Body = req.Body
	}

	if err := h.emailService.SendSync(message); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to send email",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.EmailResponse{
		Success: true,
		Message: "Email sent successfully",
	})
}

// SendBulkEmail godoc
// @Summary Send multiple emails
// @Description Send multiple emails in bulk (queued asynchronously)
// @Tags Email
// @Accept json
// @Produce json
// @Param request body models.SendBulkEmailRequest true "Bulk email details"
// @Success 200 {object} models.EmailResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/email/send-bulk [post]
// @Security BearerAuth
func (h *EmailHandler) SendBulkEmail(c *gin.Context) {
	var req models.SendBulkEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	messages := make([]email.EmailMessage, 0, len(req.Emails))
	for _, emailReq := range req.Emails {
		message := email.EmailMessage{
			To:      emailReq.To,
			Subject: emailReq.Subject,
		}

		if emailReq.IsHTML {
			message.HTMLBody = emailReq.Body
		} else {
			message.Body = emailReq.Body
		}

		messages = append(messages, message)
	}

	if err := h.emailService.SendBulk(messages); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to send bulk emails",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.EmailResponse{
		Success: true,
		Message: "Bulk emails queued successfully",
	})
}

// GetEmailStatus godoc
// @Summary Get email service status
// @Description Get the current status and queue size of the email service
// @Tags Email
// @Produce json
// @Success 200 {object} models.EmailStatusResponse
// @Router /api/v1/email/status [get]
// @Security BearerAuth
func (h *EmailHandler) GetEmailStatus(c *gin.Context) {
	queueSize := h.emailService.QueueSize()

	c.JSON(http.StatusOK, models.EmailStatusResponse{
		QueueSize: queueSize,
		Status:    "operational",
		Message:   "Email service is running",
	})
}

// TestEmailConnection godoc
// @Summary Test email service connection
// @Description Test the SMTP connection to ensure email service is working
// @Tags Email
// @Produce json
// @Success 200 {object} models.TestEmailConnectionResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/email/test-connection [get]
// @Security BearerAuth
func (h *EmailHandler) TestEmailConnection(c *gin.Context) {
	if err := h.emailService.TestConnection(); err != nil {
		c.JSON(http.StatusInternalServerError, models.TestEmailConnectionResponse{
			Success:   false,
			Message:   "Email connection failed: " + err.Error(),
			Connected: false,
		})
		return
	}

	c.JSON(http.StatusOK, models.TestEmailConnectionResponse{
		Success:   true,
		Message:   "Email connection successful",
		Connected: true,
	})
}
