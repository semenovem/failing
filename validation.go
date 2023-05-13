package failing

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

// SendValidationResponse ответ клиенту ошибкой валидации
func (s *Service) SendValidationResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, s.newValidationError(c, err))
}

// newValidationError ошибка валидации
func (s *Service) newValidationError(c echo.Context, err error) *Response {
	var (
		validationFields []*ValidationErrors
		lang             = extractLanguage(c)
	)

	if errs, ok := err.(validator.ValidationErrors); ok {
		validationFields = s.validationError(lang, errs)
		err = nil
	}

	return s.NewResponse(c, invalidRequestMessage.DefaultText, validationFields, err)
}

func (s *Service) validationError(lang msgLang, errs []validator.FieldError) []*ValidationErrors {
	errors := make([]*ValidationErrors, 0)

	for _, fieldError := range errs {
		f := ValidationErrors{Path: toSnakeCase(fieldError.StructField())}

		if msg := s.findMessageByValidationMessageTag(fieldError.Tag()); msg == nil {
			f.Message = fieldError.Translate(s.translatorDefault)
		} else {
			f.Message = msg.Text(lang)
		}

		errors = append(errors, &f)
	}

	return errors
}

func (s *Service) findMessageByValidationMessageTag(tag string) *Message {
	if msgKey, ok := s.validationMessageMap[tag]; ok {
		return s.messages[msgKey]
	}

	return nil
}
