package middleware

import (
	"net/http"

	appErrors "base-golang-restful/app/errors"
	"base-golang-restful/app/i18n"
	"base-golang-restful/app/logger"
	"base-golang-restful/app/models"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		lang := i18n.GetLanguage(c)

		var appErr *appErrors.AppError
		var localizedErr *appErrors.LocalizedError

		switch e := err.(type) {
		case *appErrors.LocalizedError:
			localizedErr = e
			appErr = localizedErr.Localize(lang)

		case *appErrors.AppError:
			appErr = e

		default:
			appErr = &appErrors.AppError{
				Code:       "INTERNAL_SERVER_ERROR",
				Message:    i18n.Translate(lang, "common.internal_error", nil),
				StatusCode: http.StatusInternalServerError,
			}
		}

		logger.Error().
			Str("correlation_id", logger.GetCorrelationID(c)).
			Str("error_code", appErr.Code).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Interface("details", appErr.Details).
			Err(appErr.Err).
			Msg(appErr.Message)

		response := models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		}

		c.JSON(appErr.StatusCode, response)
	}
}

func HandleNotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.GetLanguage(c)
		
		response := models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "NOT_FOUND",
				Message: i18n.Translate(lang, "common.not_found", nil),
			},
		}

		c.JSON(http.StatusNotFound, response)
	}
}

func HandleMethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := i18n.GetLanguage(c)
		
		response := models.APIResponse{
			Success: false,
			Error: &models.ErrorInfo{
				Code:    "METHOD_NOT_ALLOWED",
				Message: i18n.TranslateWithDefault(lang, "common.method_not_allowed", "Method not allowed", nil),
			},
		}

		c.JSON(http.StatusMethodNotAllowed, response)
	}
}
