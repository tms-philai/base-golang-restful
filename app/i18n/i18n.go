package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle       *i18n.Bundle
	localizers   map[string]*i18n.Localizer
	once         sync.Once
	localizersMu sync.RWMutex
)

type I18nConfig struct {
	DefaultLanguage string
	LocalesPath     string
	SupportedLangs  []string
}

func InitI18n(config I18nConfig) error {
	var initErr error
	once.Do(func() {
		defaultLang := config.DefaultLanguage
		if defaultLang == "" {
			defaultLang = "en"
		}

		tag, err := language.Parse(defaultLang)
		if err != nil {
			initErr = fmt.Errorf("invalid default language: %w", err)
			return
		}

		bundle = i18n.NewBundle(tag)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		localizers = make(map[string]*i18n.Localizer)

		localesPath := config.LocalesPath
		if localesPath == "" {
			localesPath = "./locales"
		}

		supportedLangs := config.SupportedLangs
		if len(supportedLangs) == 0 {
			supportedLangs = []string{"en", "vi"}
		}

		for _, lang := range supportedLangs {
			filePath := filepath.Join(localesPath, fmt.Sprintf("%s.json", lang))
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				continue
			}

			if _, err := bundle.LoadMessageFile(filePath); err != nil {
				initErr = fmt.Errorf("failed to load locale file %s: %w", filePath, err)
				return
			}

			localizers[lang] = i18n.NewLocalizer(bundle, lang)
		}

		if len(localizers) == 0 {
			initErr = fmt.Errorf("no locale files loaded from %s", localesPath)
			return
		}
	})

	return initErr
}

func GetLocalizer(lang string) *i18n.Localizer {
	if lang == "" {
		lang = "en"
	}

	localizersMu.RLock()
	localizer, exists := localizers[lang]
	localizersMu.RUnlock()

	if exists {
		return localizer
	}

	localizersMu.Lock()
	defer localizersMu.Unlock()

	localizer = i18n.NewLocalizer(bundle, lang)
	localizers[lang] = localizer

	return localizer
}

func Translate(lang, messageID string, templateData map[string]interface{}) string {
	localizer := GetLocalizer(lang)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		return messageID
	}

	return msg
}

func TranslateWithDefault(lang, messageID, defaultMsg string, templateData map[string]interface{}) string {
	localizer := GetLocalizer(lang)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    messageID,
			Other: defaultMsg,
		},
		TemplateData: templateData,
	})

	if err != nil {
		return defaultMsg
	}

	return msg
}

func MustTranslate(lang, messageID string, templateData map[string]interface{}) (string, error) {
	localizer := GetLocalizer(lang)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		return "", err
	}

	return msg, nil
}

func GetSupportedLanguages() []string {
	localizersMu.RLock()
	defer localizersMu.RUnlock()

	langs := make([]string, 0, len(localizers))
	for lang := range localizers {
		langs = append(langs, lang)
	}

	return langs
}

func IsLanguageSupported(lang string) bool {
	localizersMu.RLock()
	defer localizersMu.RUnlock()

	_, exists := localizers[lang]
	return exists
}
