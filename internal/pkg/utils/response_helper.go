package utils

import (
	"base-gin/internal/domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondSuccess(c *gin.Context, data interface{}, message string) {
	response := models.NewSuccessResponse(data, message)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusOK, response)
}

func RespondCreated(c *gin.Context, data interface{}) {
	response := models.NewCreatedResponse(data)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusCreated, response)
}

func RespondUpdated(c *gin.Context, data interface{}) {
	response := models.NewUpdatedResponse(data)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusOK, response)
}

func RespondDeleted(c *gin.Context) {
	response := models.NewDeletedResponse()
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusOK, response)
}

func RespondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func RespondList(c *gin.Context, data interface{}, pagination *models.PaginationMeta) {
	response := models.NewListResponse(data, pagination)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusOK, response)
}

func RespondError(c *gin.Context, statusCode int, code, message string, details interface{}) {
	response := models.NewErrorResponse(code, message, details)
	enrichMetadata(c, response.Metadata)
	c.JSON(statusCode, response)
}

func RespondValidationError(c *gin.Context, errors models.ValidationErrors) {
	response := models.NewValidationErrorResponse(errors)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusBadRequest, response)
}

func RespondUnauthorized(c *gin.Context) {
	response := models.NewUnauthorizedErrorResponse()
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusUnauthorized, response)
}

func RespondForbidden(c *gin.Context) {
	response := models.NewForbiddenErrorResponse()
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusForbidden, response)
}

func RespondNotFound(c *gin.Context, resource string) {
	response := models.NewNotFoundErrorResponse(resource)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusNotFound, response)
}

func RespondConflict(c *gin.Context, resource string) {
	response := models.NewConflictErrorResponse(resource)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusConflict, response)
}

func RespondInternalError(c *gin.Context) {
	response := models.NewInternalErrorResponse()
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusInternalServerError, response)
}

func RespondBadRequest(c *gin.Context, message string) {
	response := models.NewBadRequestErrorResponse(message)
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusBadRequest, response)
}

func RespondServiceUnavailable(c *gin.Context) {
	response := models.NewServiceUnavailableErrorResponse()
	enrichMetadata(c, response.Metadata)
	c.JSON(http.StatusServiceUnavailable, response)
}

func enrichMetadata(c *gin.Context, metadata *models.Metadata) {
	if metadata == nil {
		return
	}

	if correlationID, exists := c.Get("correlation_id"); exists {
		if cid, ok := correlationID.(string); ok {
			metadata.CorrelationID = cid
		}
	}

	if requestID, exists := c.Get("request_id"); exists {
		if rid, ok := requestID.(string); ok {
			metadata.RequestID = rid
		}
	}

	metadata.Version = "v1"
}
