package errors

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"base-golang-restful/app/i18n"

	"github.com/stretchr/testify/assert"
)

func setupTestI18n(t *testing.T) {
	tmpDir := t.TempDir()

	enContent := `{
  "test": {
    "message": "Test message in English",
    "with_data": "Hello, {{.Name}}!"
  }
}`

	viContent := `{
  "test": {
    "message": "Thông báo kiểm tra bằng tiếng Việt",
    "with_data": "Xin chào, {{.Name}}!"
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
}

func TestLocalizedError_Localize(t *testing.T) {
	setupTestI18n(t)

	tests := []struct {
		name            string
		localizedErr    *LocalizedError
		lang            string
		expectedMessage string
	}{
		{
			name: "localize to english",
			localizedErr: &LocalizedError{
				AppError: &AppError{
					Code:       "TEST_ERROR",
					StatusCode: http.StatusBadRequest,
				},
				MessageKey: "test.message",
			},
			lang:            "en",
			expectedMessage: "Test message in English",
		},
		{
			name: "localize to vietnamese",
			localizedErr: &LocalizedError{
				AppError: &AppError{
					Code:       "TEST_ERROR",
					StatusCode: http.StatusBadRequest,
				},
				MessageKey: "test.message",
			},
			lang:            "vi",
			expectedMessage: "Thông báo kiểm tra bằng tiếng Việt",
		},
		{
			name: "localize with template data",
			localizedErr: &LocalizedError{
				AppError: &AppError{
					Code:       "TEST_ERROR",
					StatusCode: http.StatusBadRequest,
				},
				MessageKey: "test.with_data",
				TemplateData: map[string]interface{}{
					"Name": "John",
				},
			},
			lang:            "en",
			expectedMessage: "Hello, John!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := tt.localizedErr.Localize(tt.lang)

			assert.Equal(t, tt.localizedErr.Code, appErr.Code)
			assert.Equal(t, tt.expectedMessage, appErr.Message)
			assert.Equal(t, tt.localizedErr.StatusCode, appErr.StatusCode)
		})
	}
}

func TestNewLocalizedError(t *testing.T) {
	code := "TEST_ERROR"
	messageKey := "test.message"
	statusCode := http.StatusBadRequest
	templateData := map[string]interface{}{"key": "value"}

	localizedErr := NewLocalizedError(code, messageKey, statusCode, templateData)

	assert.Equal(t, code, localizedErr.Code)
	assert.Equal(t, messageKey, localizedErr.MessageKey)
	assert.Equal(t, statusCode, localizedErr.StatusCode)
	assert.Equal(t, templateData, localizedErr.TemplateData)
}

func TestLocalizedErrorConstructors(t *testing.T) {
	tests := []struct {
		name           string
		constructor    func(string, map[string]interface{}) *LocalizedError
		messageKey     string
		expectedCode   string
		expectedStatus int
	}{
		{
			name:           "LocalizedBadRequest",
			constructor:    LocalizedBadRequest,
			messageKey:     "test.message",
			expectedCode:   "BAD_REQUEST",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "LocalizedUnauthorized",
			constructor:    LocalizedUnauthorized,
			messageKey:     "test.message",
			expectedCode:   "UNAUTHORIZED",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "LocalizedForbidden",
			constructor:    LocalizedForbidden,
			messageKey:     "test.message",
			expectedCode:   "FORBIDDEN",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "LocalizedNotFound",
			constructor:    LocalizedNotFound,
			messageKey:     "test.message",
			expectedCode:   "NOT_FOUND",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "LocalizedValidation",
			constructor:    LocalizedValidation,
			messageKey:     "test.message",
			expectedCode:   "VALIDATION_ERROR",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "LocalizedInternalServer",
			constructor:    LocalizedInternalServer,
			messageKey:     "test.message",
			expectedCode:   "INTERNAL_SERVER_ERROR",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateData := map[string]interface{}{"key": "value"}
			err := tt.constructor(tt.messageKey, templateData)

			assert.Equal(t, tt.expectedCode, err.Code)
			assert.Equal(t, tt.messageKey, err.MessageKey)
			assert.Equal(t, tt.expectedStatus, err.StatusCode)
			assert.Equal(t, templateData, err.TemplateData)
		})
	}
}
