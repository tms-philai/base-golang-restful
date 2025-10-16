package unit

import (
	"base-gin/internal/app/handlers"
	"base-gin/internal/domain/models"
	"base-gin/internal/pkg/utils"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock implementation of EmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendAsync(message utils.EmailMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockEmailService) SendSync(message utils.EmailMessage) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockEmailService) SendBulk(messages []utils.EmailMessage) error {
	args := m.Called(messages)
	return args.Error(0)
}

func (m *MockEmailService) QueueSize() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockEmailService) TestConnection() error {
	args := m.Called()
	return args.Error(0)
}

func TestEmailHandler_SendEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.SendEmailRequest
		mockSetup      func(*MockEmailService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful send email",
			request: models.SendEmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Body:    "Test message body",
				IsHTML:  false,
			},
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("SendAsync", mock.AnythingOfType("utils.EmailMessage")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful send HTML email",
			request: models.SendEmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Body:    "<h1>Test HTML message</h1>",
				IsHTML:  true,
			},
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("SendAsync", mock.AnythingOfType("utils.EmailMessage")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "validation error - missing recipients",
			request: models.SendEmailRequest{
				Subject: "Test Subject",
				Body:    "Test message body",
				IsHTML:  false,
			},
			mockSetup:      func(*MockEmailService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
		{
			name: "email service error",
			request: models.SendEmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Body:    "Test message body",
				IsHTML:  false,
			},
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("SendAsync", mock.AnythingOfType("utils.EmailMessage")).Return(errors.New("SMTP connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to send email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockEmailService := new(MockEmailService)
			tt.mockSetup(mockEmailService)

			// Create handler
			handler := handlers.NewEmailHandler(mockEmailService)

			// Setup router
			router := gin.New()
			router.POST("/email/send", handler.SendEmail)

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/email/send", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockEmailService.AssertExpectations(t)
		})
	}
}

func TestEmailHandler_SendEmailSync(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.SendEmailRequest
		mockSetup      func(*MockEmailService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful send sync email",
			request: models.SendEmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Body:    "Test message body",
				IsHTML:  false,
			},
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("SendSync", mock.AnythingOfType("utils.EmailMessage")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "email service error",
			request: models.SendEmailRequest{
				To:      []string{"test@example.com"},
				Subject: "Test Subject",
				Body:    "Test message body",
				IsHTML:  false,
			},
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("SendSync", mock.AnythingOfType("utils.EmailMessage")).Return(errors.New("SMTP connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to send email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockEmailService := new(MockEmailService)
			tt.mockSetup(mockEmailService)

			// Create handler
			handler := handlers.NewEmailHandler(mockEmailService)

			// Setup router
			router := gin.New()
			router.POST("/email/send-sync", handler.SendEmailSync)

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/email/send-sync", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockEmailService.AssertExpectations(t)
		})
	}
}

func TestEmailHandler_SendBulkEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        models.SendBulkEmailRequest
		mockSetup      func(*MockEmailService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful send bulk email",
			request: models.SendBulkEmailRequest{
				Emails: []models.SendEmailRequest{
					{
						To:      []string{"user1@example.com"},
						Subject: "Test Subject 1",
						Body:    "Test message 1",
						IsHTML:  false,
					},
					{
						To:      []string{"user2@example.com"},
						Subject: "Test Subject 2",
						Body:    "Test message 2",
						IsHTML:  false,
					},
				},
			},
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("SendBulk", mock.AnythingOfType("[]utils.EmailMessage")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "validation error - empty emails list",
			request: models.SendBulkEmailRequest{
				Emails: []models.SendEmailRequest{},
			},
			mockSetup:      func(*MockEmailService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request",
		},
		{
			name: "email service error",
			request: models.SendBulkEmailRequest{
				Emails: []models.SendEmailRequest{
					{
						To:      []string{"user1@example.com"},
						Subject: "Test Subject 1",
						Body:    "Test message 1",
						IsHTML:  false,
					},
				},
			},
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("SendBulk", mock.AnythingOfType("[]utils.EmailMessage")).Return(errors.New("SMTP connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to send bulk emails",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockEmailService := new(MockEmailService)
			tt.mockSetup(mockEmailService)

			// Create handler
			handler := handlers.NewEmailHandler(mockEmailService)

			// Setup router
			router := gin.New()
			router.POST("/email/send-bulk", handler.SendBulkEmail)

			// Create request
			jsonBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/email/send-bulk", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response models.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, response.Error)
			}

			// Verify mock expectations
			mockEmailService.AssertExpectations(t)
		})
	}
}

func TestEmailHandler_GetEmailStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		mockSetup         func(*MockEmailService)
		expectedStatus    int
		expectedQueueSize int
	}{
		{
			name: "successful get email status",
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("QueueSize").Return(5)
			},
			expectedStatus:    http.StatusOK,
			expectedQueueSize: 5,
		},
		{
			name: "empty queue",
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("QueueSize").Return(0)
			},
			expectedStatus:    http.StatusOK,
			expectedQueueSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockEmailService := new(MockEmailService)
			tt.mockSetup(mockEmailService)

			// Create handler
			handler := handlers.NewEmailHandler(mockEmailService)

			// Setup router
			router := gin.New()
			router.GET("/email/status", handler.GetEmailStatus)

			// Create request
			req, _ := http.NewRequest("GET", "/email/status", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.EmailStatusResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedQueueSize, response.QueueSize)
			assert.Equal(t, "operational", response.Status)

			// Verify mock expectations
			mockEmailService.AssertExpectations(t)
		})
	}
}

func TestEmailHandler_TestEmailConnection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		mockSetup         func(*MockEmailService)
		expectedStatus    int
		expectedSuccess   bool
		expectedConnected bool
	}{
		{
			name: "successful connection test",
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("TestConnection").Return(nil)
			},
			expectedStatus:    http.StatusOK,
			expectedSuccess:   true,
			expectedConnected: true,
		},
		{
			name: "connection test failed",
			mockSetup: func(emailService *MockEmailService) {
				emailService.On("TestConnection").Return(errors.New("SMTP connection failed"))
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedSuccess:   false,
			expectedConnected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockEmailService := new(MockEmailService)
			tt.mockSetup(mockEmailService)

			// Create handler
			handler := handlers.NewEmailHandler(mockEmailService)

			// Setup router
			router := gin.New()
			router.GET("/email/test-connection", handler.TestEmailConnection)

			// Create request
			req, _ := http.NewRequest("GET", "/email/test-connection", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response models.TestEmailConnectionResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSuccess, response.Success)
			assert.Equal(t, tt.expectedConnected, response.Connected)

			// Verify mock expectations
			mockEmailService.AssertExpectations(t)
		})
	}
}
