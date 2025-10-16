package errors

import "base-gin/internal/pkg/i18n"

type LocalizedError struct {
	*AppError
	MessageKey   string                 `json:"-"`
	TemplateData map[string]interface{} `json:"-"`
}

func (e *LocalizedError) Localize(lang string) *AppError {
	localizedMessage := i18n.Translate(lang, e.MessageKey, e.TemplateData)

	return &AppError{
		Code:       e.Code,
		Message:    localizedMessage,
		StatusCode: e.StatusCode,
		Details:    e.Details,
		Err:        e.Err,
	}
}

func NewLocalizedError(code, messageKey string, statusCode int, templateData map[string]interface{}) *LocalizedError {
	return &LocalizedError{
		AppError: &AppError{
			Code:       code,
			StatusCode: statusCode,
		},
		MessageKey:   messageKey,
		TemplateData: templateData,
	}
}

func LocalizedBadRequest(messageKey string, templateData map[string]interface{}) *LocalizedError {
	return &LocalizedError{
		AppError: &AppError{
			Code:       ErrBadRequest.Code,
			StatusCode: ErrBadRequest.StatusCode,
		},
		MessageKey:   messageKey,
		TemplateData: templateData,
	}
}

func LocalizedUnauthorized(messageKey string, templateData map[string]interface{}) *LocalizedError {
	return &LocalizedError{
		AppError: &AppError{
			Code:       ErrUnauthorized.Code,
			StatusCode: ErrUnauthorized.StatusCode,
		},
		MessageKey:   messageKey,
		TemplateData: templateData,
	}
}

func LocalizedForbidden(messageKey string, templateData map[string]interface{}) *LocalizedError {
	return &LocalizedError{
		AppError: &AppError{
			Code:       ErrForbidden.Code,
			StatusCode: ErrForbidden.StatusCode,
		},
		MessageKey:   messageKey,
		TemplateData: templateData,
	}
}

func LocalizedNotFound(messageKey string, templateData map[string]interface{}) *LocalizedError {
	return &LocalizedError{
		AppError: &AppError{
			Code:       ErrNotFound.Code,
			StatusCode: ErrNotFound.StatusCode,
		},
		MessageKey:   messageKey,
		TemplateData: templateData,
	}
}

func LocalizedValidation(messageKey string, templateData map[string]interface{}) *LocalizedError {
	return &LocalizedError{
		AppError: &AppError{
			Code:       ErrValidation.Code,
			StatusCode: ErrValidation.StatusCode,
		},
		MessageKey:   messageKey,
		TemplateData: templateData,
	}
}

func LocalizedInternalServer(messageKey string, templateData map[string]interface{}) *LocalizedError {
	return &LocalizedError{
		AppError: &AppError{
			Code:       ErrInternalServer.Code,
			StatusCode: ErrInternalServer.StatusCode,
		},
		MessageKey:   messageKey,
		TemplateData: templateData,
	}
}
