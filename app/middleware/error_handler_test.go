package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	appErrors "base-golang-restful-app/errors"
	"base-golang-restful-app/i18n"
	"base-golang-restful-app/logger"
	"base-golang-restful-app/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestEnvironment(t *testing.T) {
	tmpDir := t.TempDir()

	enContent := `{
  "common": {
    "internal_error": "Internal server error",
    "not_found": "Resource not found",
    "method_not_allowed": "Method not allowed"
  }
}`

	viContent := `{
  "common": {
    "internal_error": "Lỗi máy chủ nội bộ",
    "not_found": "Không tìm thấy tài nguyên",
    "method_not_allowed": "Phương thức không được phép"
  }
}`

	err := os.WriteFile(filepath.Join(tmpDir, "en.json"), []byte(enContent), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "vi.json"), []byte(viContent), 0644)
	assert.NoError(t, err)

	err = i18n.InitI18n(i18n.I18nConfig{
		DefaultLanguage: "en",
		LocalesPath:     tmpDir,
		SupportedLangs:  []string{"en", "vi"},
	})
	assert.NoError(t, err)

	logDir := t.TempDir()
	err = logger.InitLogger(logger.LogConfig{
		Level:      "info",
		OutputPath: filepath.Join(logDir, "test.log"),
	})
	assert.NoError(t, err)
}

func TestErrorHandler_AppError(t *testing.T) {
	setupTestEnvironment(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(appErrors.NotFound("User not found"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "NOT_FOUND", response.Error.Code)
	assert.Equal(t, "User not found", response.Error.Message)
}

func TestErrorHandler_LocalizedError(t *testing.T) {
	setupTestEnvironment(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(appErrors.LocalizedNotFound("common.not_found", nil))
	})

	tests := []struct {
		name            string
		lang            string
		expectedMessage string
	}{
		{
			name:            "english",
			lang:            "en",
			expectedMessage: "Resource not found",
		},
		{
			name:            "vietnamese",
			lang:            "vi",
			expectedMessage: "Không tìm thấy tài nguyên",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test?lang="+tt.lang, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)

			var response models.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.False(t, response.Success)
			assert.NotNil(t, response.Error)
			assert.Equal(t, "NOT_FOUND", response.Error.Code)
			assert.Equal(t, tt.expectedMessage, response.Error.Message)
		})
	}
}

func TestErrorHandler_GenericError(t *testing.T) {
	setupTestEnvironment(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New("some generic error"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "INTERNAL_SERVER_ERROR", response.Error.Code)
}

func TestHandleNotFound(t *testing.T) {
	setupTestEnvironment(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.NoRoute(HandleNotFound())

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "NOT_FOUND", response.Error.Code)
}

func TestHandleMethodNotAllowed(t *testing.T) {
	setupTestEnvironment(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.NoMethod(HandleMethodNotAllowed())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "METHOD_NOT_ALLOWED", response.Error.Code)
}
