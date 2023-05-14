package failing

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

// SendValidationResponse ответ клиенту ошибкой валидации
func (s *Service) SendValidationResponse(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, s.newValidationResponse(c, err))
}

func (s *Service) newValidationResponse(c echo.Context, err error) *Response {
	var (
		validationFields []*ValidationError
		lang             = extractLanguage(c)
	)

	if errs, ok := err.(validator.ValidationErrors); ok {
		validationFields = s.validationErrors(lang, errs)
		err = nil
	}

	resp := s.NewResponse(c, invalidRequestMessage.DefaultText, validationFields, err)
	resp.ValidationErrors = validationFields

	return resp
}

func (s *Service) validationErrors(lang msgLang, errs []validator.FieldError) []*ValidationError {
	validationErrs := make([]*ValidationError, 0)

	for _, fieldError := range errs {
		f := ValidationError{Path: toSnakeCase(fieldError.StructField())}

		if msg := s.findMessageByValidationMessageTag(fieldError.Tag()); msg == nil {
			f.Message = fieldError.Translate(s.translatorDefault)
		} else {
			f.Message = msg.Text(lang)
		}

		validationErrs = append(validationErrs, &f)
	}

	return validationErrs
}

func (s *Service) findMessageByValidationMessageTag(tag string) *Message {
	if msgKey, ok := s.validationMessageMap[tag]; ok {
		return s.messages[msgKey]
	}

	return nil
}
