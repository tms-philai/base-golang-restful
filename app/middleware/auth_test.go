package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"base-golang-restful-app/auth"
	"base-golang-restful-app/i18n"
	"base-golang-restful-app/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupAuthTestEnvironment(t *testing.T) *auth.JWTManager {
	tmpDir := t.TempDir()

	enContent := `{
  "auth": {
    "token_expired": "Token has expired",
    "token_invalid": "Invalid token"
  }
}`

	err := os.WriteFile(filepath.Join(tmpDir, "en.json"), []byte(enContent), 0644)
	assert.NoError(t, err)

	err = i18n.InitI18n(i18n.I18nConfig{
		DefaultLanguage: "en",
		LocalesPath:     tmpDir,
		SupportedLangs:  []string{"en"},
	})
	assert.NoError(t, err)

	logDir := t.TempDir()
	err = logger.InitLogger(logger.LogConfig{
		Level:      "info",
		OutputPath: filepath.Join(logDir, "test.log"),
	})
	assert.NoError(t, err)

	return auth.NewJWTManager(auth.JWTConfig{
		SecretKey:           "test-secret",
		AccessTokenDuration: 1 * time.Hour,
	})
}

func TestAuthMiddleware_Authenticate_ValidToken(t *testing.T) {
	jwtManager := setupAuthTestEnvironment(t)
	middleware := NewAuthMiddleware(jwtManager)

	userID := uuid.New()
	email := "test@example.com"
	token, err := jwtManager.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.Use(middleware.Authenticate())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_Authenticate_NoToken(t *testing.T) {
	jwtManager := setupAuthTestEnvironment(t)
	middleware := NewAuthMiddleware(jwtManager)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.Use(ErrorHandler())
	router.Use(middleware.Authenticate())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Authenticate_InvalidToken(t *testing.T) {
	jwtManager := setupAuthTestEnvironment(t)
	middleware := NewAuthMiddleware(jwtManager)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.Use(ErrorHandler())
	router.Use(middleware.Authenticate())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_Authenticate_ExpiredToken(t *testing.T) {
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		SecretKey:           "test-secret",
		AccessTokenDuration: 1 * time.Millisecond,
	})

	setupAuthTestEnvironment(t)
	middleware := NewAuthMiddleware(jwtManager)

	userID := uuid.New()
	email := "test@example.com"
	token, err := jwtManager.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(i18n.LanguageMiddleware())
	router.Use(ErrorHandler())
	router.Use(middleware.Authenticate())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_OptionalAuthenticate(t *testing.T) {
	jwtManager := setupAuthTestEnvironment(t)
	middleware := NewAuthMiddleware(jwtManager)

	userID := uuid.New()
	email := "test@example.com"
	token, err := jwtManager.GenerateAccessToken(userID, email)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		shouldHaveUser bool
	}{
		{
			name:           "with valid token",
			token:          token,
			expectedStatus: http.StatusOK,
			shouldHaveUser: true,
		},
		{
			name:           "without token",
			token:          "",
			expectedStatus: http.StatusOK,
			shouldHaveUser: false,
		},
		{
			name:           "with invalid token",
			token:          "invalid.token",
			expectedStatus: http.StatusOK,
			shouldHaveUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(middleware.OptionalAuthenticate())
			router.GET("/test", func(c *gin.Context) {
				_, hasUser := c.Get("user_id")
				c.JSON(http.StatusOK, gin.H{"has_user": hasUser})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestExtractToken(t *testing.T) {
	tests := []struct {
		name          string
		setupRequest  func(*http.Request)
		expectedToken string
	}{
		{
			name: "from authorization header",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer test-token")
			},
			expectedToken: "test-token",
		},
		{
			name: "from query parameter",
			setupRequest: func(req *http.Request) {
				q := req.URL.Query()
				q.Add("token", "test-token")
				req.URL.RawQuery = q.Encode()
			},
			expectedToken: "test-token",
		},
		{
			name: "from cookie",
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  "access_token",
					Value: "test-token",
				})
			},
			expectedToken: "test-token",
		},
		{
			name:          "no token",
			setupRequest:  func(req *http.Request) {},
			expectedToken: "",
		},
		{
			name: "invalid authorization format",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "InvalidFormat")
			},
			expectedToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupRequest(req)
			c.Request = req

			token := extractToken(c)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	userID := uuid.New().String()
	c.Set("user_id", userID)

	retrievedID, exists := GetUserID(c)
	assert.True(t, exists)
	assert.Equal(t, userID, retrievedID)

	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, exists = GetUserID(c2)
	assert.False(t, exists)
}

func TestGetUserEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	email := "test@example.com"
	c.Set("user_email", email)

	retrievedEmail, exists := GetUserEmail(c)
	assert.True(t, exists)
	assert.Equal(t, email, retrievedEmail)

	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, exists = GetUserEmail(c2)
	assert.False(t, exists)
}

func TestGetTokenClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	claims := &auth.JWTClaims{
		UserID: uuid.New(),
		Email:  "test@example.com",
	}
	c.Set("token_claims", claims)

	retrievedClaims, exists := GetTokenClaims(c)
	assert.True(t, exists)
	assert.Equal(t, claims, retrievedClaims)

	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, exists = GetTokenClaims(c2)
	assert.False(t, exists)
}
