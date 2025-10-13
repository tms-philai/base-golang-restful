package handlers

import (
	"net/http"
	"time"

	"base-golang-restful-app/email"
	"base-golang-restful-app/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService *email.NotificationService
	inAppChannel        *email.InAppNotificationChannel
}

func NewNotificationHandler(notificationService *email.NotificationService, inAppChannel *email.InAppNotificationChannel) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		inAppChannel:        inAppChannel,
	}
}

// SendNotification godoc
// @Summary Send a notification
// @Description Send a notification to a user via specified channel (email or in-app)
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body models.SendNotificationRequest true "Notification details"
// @Success 200 {object} models.NotificationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/notifications/send [post]
// @Security BearerAuth
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req models.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid UUID",
		})
		return
	}

	if err := h.notificationService.NotifyUser(
		c.Request.Context(),
		userID,
		email.NotificationType(req.Type),
		req.Title,
		req.Message,
		req.Data,
	); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to send notification",
			Message: err.Error(),
		})
		return
	}

	now := time.Now()
	c.JSON(http.StatusOK, models.NotificationResponse{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      req.Type,
		Status:    models.NotificationStatusSent,
		Title:     req.Title,
		Message:   req.Message,
		SentAt:    &now,
		CreatedAt: now,
	})
}

// SendBulkNotification godoc
// @Summary Send multiple notifications
// @Description Send multiple notifications in bulk
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body models.SendBulkNotificationRequest true "Bulk notification details"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/notifications/send-bulk [post]
// @Security BearerAuth
func (h *NotificationHandler) SendBulkNotification(c *gin.Context) {
	var req models.SendBulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	for _, notif := range req.Notifications {
		userID, err := uuid.Parse(notif.UserID)
		if err != nil {
			continue
		}

		_ = h.notificationService.NotifyUser(
			c.Request.Context(),
			userID,
			email.NotificationType(notif.Type),
			notif.Title,
			notif.Message,
			notif.Data,
		)
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Bulk notifications sent successfully",
	})
}

// GetNotifications godoc
// @Summary Get user notifications
// @Description Get all notifications for a specific user
// @Tags Notifications
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} models.ListNotificationsResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/notifications/user/{user_id} [get]
// @Security BearerAuth
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid UUID",
		})
		return
	}

	notifications := h.inAppChannel.GetNotifications(userID)

	response := models.ListNotificationsResponse{
		Notifications: make([]models.NotificationResponse, 0),
		Total:         len(notifications),
		Unread:        0,
	}

	for _, notif := range notifications {
		if notif.ReadAt == nil {
			response.Unread++
		}

		response.Notifications = append(response.Notifications, models.NotificationResponse{
			ID:        notif.ID,
			UserID:    notif.UserID,
			Type:      models.NotificationType(notif.Type),
			Status:    models.NotificationStatus(notif.Status),
			Title:     notif.Title,
			Message:   notif.Message,
			SentAt:    notif.SentAt,
			ReadAt:    notif.ReadAt,
			CreatedAt: notif.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a specific notification as read
// @Tags Notifications
// @Accept json
// @Produce json
// @Param request body models.MarkAsReadRequest true "Notification ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/v1/notifications/mark-read [post]
// @Security BearerAuth
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	var req models.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	notificationID, err := uuid.Parse(req.NotificationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid notification ID",
			Message: "Notification ID must be a valid UUID",
		})
		return
	}

	if err := h.inAppChannel.MarkAsRead(notificationID); err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Notification not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Notification marked as read",
	})
}

// GetServiceStatus godoc
// @Summary Get notification service status
// @Description Get the current status of notification service and available channels
// @Tags Notifications
// @Produce json
// @Success 200 {object} models.NotificationServiceStatusResponse
// @Router /api/v1/notifications/status [get]
// @Security BearerAuth
func (h *NotificationHandler) GetServiceStatus(c *gin.Context) {
	c.JSON(http.StatusOK, models.NotificationServiceStatusResponse{
		EmailChannel: true,
		InAppChannel: true,
		Status:       "operational",
		Message:      "Notification service is running",
	})
}
