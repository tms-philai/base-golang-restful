package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestLanguageMiddleware(t *testing.T) {
	tmpDir := setupTestLocales(t)

	bundle = nil
	localizers = nil
	once = sync.Once{}

	err := InitI18n(I18nConfig{
		DefaultLanguage: "en",
		LocalesPath:     tmpDir,
		SupportedLangs:  []string{"en", "vi"},
	})
	assert.NoError(t, err)

	tests := []struct {
		name           string
		queryParam     string
		header         string
		cookie         string
		expectedLang   string
	}{
		{
			name:         "query parameter takes precedence",
			queryParam:   "vi",
			header:       "en",
			expectedLang: "vi",
		},
		{
			name:         "header when no query param",
			header:       "vi",
			expectedLang: "vi",
		},
		{
			name:         "cookie when no query or header",
			cookie:       "vi",
			expectedLang: "vi",
		},
		{
			name:         "default when nothing provided",
			expectedLang: "en",
		},
		{
			name:         "unsupported language falls back to default",
			queryParam:   "fr",
			expectedLang: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.Use(LanguageMiddleware())
			router.GET("/test", func(c *gin.Context) {
				lang := GetLanguage(c)
				c.String(http.StatusOK, lang)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Add("lang", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}
			if tt.header != "" {
				req.Header.Set("Accept-Language", tt.header)
			}
			if tt.cookie != "" {
				req.AddCookie(&http.Cookie{
					Name:  "language",
					Value: tt.cookie,
				})
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, tt.expectedLang, w.Body.String())
		})
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		name       string
		acceptLang string
		expected   string
	}{
		{
			name:       "simple language",
			acceptLang: "en",
			expected:   "en",
		},
		{
			name:       "language with region",
			acceptLang: "en-US",
			expected:   "en",
		},
		{
			name:       "language with quality",
			acceptLang: "en;q=0.9",
			expected:   "en",
		},
		{
			name:       "multiple languages",
			acceptLang: "vi,en-US;q=0.9,en;q=0.8",
			expected:   "vi",
		},
		{
			name:       "empty string",
			acceptLang: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAcceptLanguage(tt.acceptLang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetLanguage(t *testing.T) {
	tests := []struct {
		name     string
		setLang  bool
		langVal  interface{}
		expected string
	}{
		{
			name:     "valid language",
			setLang:  true,
			langVal:  "vi",
			expected: "vi",
		},
		{
			name:     "no language set",
			setLang:  false,
			expected: "en",
		},
		{
			name:     "invalid type",
			setLang:  true,
			langVal:  123,
			expected: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())

			if tt.setLang {
				c.Set(LanguageContextKey, tt.langVal)
			}

			result := GetLanguage(c)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestT(t *testing.T) {
	tmpDir := setupTestLocales(t)

	bundle = nil
	localizers = nil
	once = sync.Once{}

	err := InitI18n(I18nConfig{
		DefaultLanguage: "en",
		LocalesPath:     tmpDir,
		SupportedLangs:  []string{"en", "vi"},
	})
	assert.NoError(t, err)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(LanguageContextKey, "en")

	result := T(c, "test.hello", nil)
	assert.Equal(t, "Hello", result)

	c.Set(LanguageContextKey, "vi")
	result = T(c, "test.hello", nil)
	assert.Equal(t, "Xin chào", result)
}

func TestTWithDefault(t *testing.T) {
	tmpDir := setupTestLocales(t)

	bundle = nil
	localizers = nil
	once = sync.Once{}

	err := InitI18n(I18nConfig{
		DefaultLanguage: "en",
		LocalesPath:     tmpDir,
		SupportedLangs:  []string{"en", "vi"},
	})
	assert.NoError(t, err)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set(LanguageContextKey, "en")

	result := TWithDefault(c, "test.hello", "Default", nil)
	assert.Equal(t, "Hello", result)

	result = TWithDefault(c, "test.missing", "Default Message", nil)
	assert.Equal(t, "Default Message", result)
}
