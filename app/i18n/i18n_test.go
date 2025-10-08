package i18n

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestLocales(t *testing.T) string {
	tmpDir := t.TempDir()

	enContent := `{
  "test": {
    "hello": "Hello",
    "greeting": "Hello, {{.Name}}!",
    "count": "You have {{.Count}} messages"
  }
}`

	viContent := `{
  "test": {
    "hello": "Xin chào",
    "greeting": "Xin chào, {{.Name}}!",
    "count": "Bạn có {{.Count}} tin nhắn"
  }
}`

	err := os.WriteFile(filepath.Join(tmpDir, "en.json"), []byte(enContent), 0644)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "vi.json"), []byte(viContent), 0644)
	assert.NoError(t, err)

	return tmpDir
}

func TestInitI18n(t *testing.T) {
	tmpDir := setupTestLocales(t)

	tests := []struct {
		name        string
		config      I18nConfig
		expectError bool
	}{
		{
			name: "valid config",
			config: I18nConfig{
				DefaultLanguage: "en",
				LocalesPath:     tmpDir,
				SupportedLangs:  []string{"en", "vi"},
			},
			expectError: false,
		},
		{
			name: "default values",
			config: I18nConfig{
				LocalesPath:    tmpDir,
				SupportedLangs: []string{"en", "vi"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle = nil
			localizers = nil
			once = sync.Once{}

			err := InitI18n(tt.config)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bundle)
				assert.NotNil(t, localizers)
			}
		})
	}
}

func TestGetLocalizer(t *testing.T) {
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
		name     string
		lang     string
		expected bool
	}{
		{
			name:     "english localizer",
			lang:     "en",
			expected: true,
		},
		{
			name:     "vietnamese localizer",
			lang:     "vi",
			expected: true,
		},
		{
			name:     "default for empty",
			lang:     "",
			expected: true,
		},
		{
			name:     "unsupported language",
			lang:     "fr",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localizer := GetLocalizer(tt.lang)
			assert.NotNil(t, localizer)
		})
	}
}

func TestTranslate(t *testing.T) {
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
		name         string
		lang         string
		messageID    string
		templateData map[string]interface{}
		expected     string
	}{
		{
			name:         "simple english translation",
			lang:         "en",
			messageID:    "test.hello",
			templateData: nil,
			expected:     "Hello",
		},
		{
			name:         "simple vietnamese translation",
			lang:         "vi",
			messageID:    "test.hello",
			templateData: nil,
			expected:     "Xin chào",
		},
		{
			name:      "english with template",
			lang:      "en",
			messageID: "test.greeting",
			templateData: map[string]interface{}{
				"Name": "John",
			},
			expected: "Hello, John!",
		},
		{
			name:      "vietnamese with template",
			lang:      "vi",
			messageID: "test.greeting",
			templateData: map[string]interface{}{
				"Name": "Minh",
			},
			expected: "Xin chào, Minh!",
		},
		{
			name:      "missing message returns ID",
			lang:      "en",
			messageID: "test.nonexistent",
			expected:  "test.nonexistent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Translate(tt.lang, tt.messageID, tt.templateData)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTranslateWithDefault(t *testing.T) {
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
		name         string
		lang         string
		messageID    string
		defaultMsg   string
		templateData map[string]interface{}
		expected     string
	}{
		{
			name:       "existing message",
			lang:       "en",
			messageID:  "test.hello",
			defaultMsg: "Default Hello",
			expected:   "Hello",
		},
		{
			name:       "missing message uses default",
			lang:       "en",
			messageID:  "test.missing",
			defaultMsg: "Default Message",
			expected:   "Default Message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TranslateWithDefault(tt.lang, tt.messageID, tt.defaultMsg, tt.templateData)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
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

	langs := GetSupportedLanguages()
	assert.Len(t, langs, 2)
	assert.Contains(t, langs, "en")
	assert.Contains(t, langs, "vi")
}

func TestIsLanguageSupported(t *testing.T) {
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
		name     string
		lang     string
		expected bool
	}{
		{
			name:     "supported english",
			lang:     "en",
			expected: true,
		},
		{
			name:     "supported vietnamese",
			lang:     "vi",
			expected: true,
		},
		{
			name:     "unsupported french",
			lang:     "fr",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLanguageSupported(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}
