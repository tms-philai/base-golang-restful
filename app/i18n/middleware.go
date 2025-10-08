package i18n

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	LanguageContextKey = "language"
	DefaultLanguage    = "en"
)

func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := extractLanguage(c)

		if !IsLanguageSupported(lang) {
			lang = DefaultLanguage
		}

		c.Set(LanguageContextKey, lang)

		c.Next()
	}
}

func extractLanguage(c *gin.Context) string {
	lang := c.Query("lang")
	if lang != "" {
		return lang
	}

	lang = c.GetHeader("Accept-Language")
	if lang != "" {
		lang = parseAcceptLanguage(lang)
		if lang != "" {
			return lang
		}
	}

	cookie, err := c.Cookie("language")
	if err == nil && cookie != "" {
		return cookie
	}

	return DefaultLanguage
}

func parseAcceptLanguage(acceptLang string) string {
	if acceptLang == "" {
		return ""
	}

	parts := strings.Split(acceptLang, ",")
	if len(parts) == 0 {
		return ""
	}

	firstLang := strings.TrimSpace(parts[0])

	langParts := strings.Split(firstLang, ";")
	if len(langParts) == 0 {
		return ""
	}

	lang := strings.TrimSpace(langParts[0])

	if strings.Contains(lang, "-") {
		lang = strings.Split(lang, "-")[0]
	}

	return lang
}

func GetLanguage(c *gin.Context) string {
	lang, exists := c.Get(LanguageContextKey)
	if !exists {
		return DefaultLanguage
	}

	langStr, ok := lang.(string)
	if !ok {
		return DefaultLanguage
	}

	return langStr
}

func T(c *gin.Context, messageID string, templateData map[string]interface{}) string {
	lang := GetLanguage(c)
	return Translate(lang, messageID, templateData)
}

func TWithDefault(c *gin.Context, messageID, defaultMsg string, templateData map[string]interface{}) string {
	lang := GetLanguage(c)
	return TranslateWithDefault(lang, messageID, defaultMsg, templateData)
}
