package failing

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Config struct {
	IsDevMode            bool // Режим разработки
	TranslatorDefault    ut.Translator
	Translators          map[msgLang]ut.Translator
	Messages             map[string]*Message
	ValidationMessageMap map[string]string
	Logger               logger
	HTTPStatuses         map[int]*Message
}

type Service struct {
	translatorDefault    ut.Translator
	translators          map[msgLang]ut.Translator
	isDev                bool
	messages             map[string]*Message
	validationMessageMap map[string]string
	logger               logger
	httpStatuses         map[int]*Message
}

func New(c *Config) *Service {
	s := &Service{
		translatorDefault:    c.TranslatorDefault,
		translators:          c.Translators,
		isDev:                c.IsDevMode,
		messages:             c.Messages,
		validationMessageMap: c.ValidationMessageMap,
	}

	if c.HTTPStatuses == nil {
		s.httpStatuses = statuses
	}

	return s
}

// NewResponse
// opts, параметры, определяемые по типу:
// > map[string]interface{} - дополнительные поля
// > error - ошибка, если isDev = true (dev режим) добавить информацию в доп поля
// > []ValidationError - ошибки валидации
func (s *Service) NewResponse(c echo.Context, opts ...interface{}) *Response {
	return s.newResponse(extractLanguage(c), s.parseOpts(opts))
}

func (s *Service) SendResponse(c echo.Context, httpStatus int, opts ...interface{}) error {
	var (
		lang = extractLanguage(c)
		opt  = s.parseOpts(opts)
	)

	if opt.message == nil {
		if st, ok := s.httpStatuses[httpStatus]; ok {
			opt.message = st
		} else {
			opt.message = unknownMessage
		}
	}

	return c.JSON(httpStatus, s.newResponse(lang, opt))
}

func (s *Service) SendStatusInternalServerResponse(c echo.Context, err error) error {
	return s.SendResponse(c, http.StatusInternalServerError, err)
}

// SendFromNestedResponse обработка ошибки из вложенных вызовов
func (s *Service) SendFromNestedResponse(c echo.Context, nestedResp *NestedResponse) error {
	var (
		lang = extractLanguage(c)
		opt  = s.parseOpts(nestedResp.opts)
	)

	if opt.message == nil {
		if st, ok := s.httpStatuses[nestedResp.httpStatusCode]; ok {
			opt.message = st
		} else {
			opt.message = unknownMessage
		}
	}

	response := s.newResponse(lang, opt)

	if errs, ok := nestedResp.validationError.(validator.ValidationErrors); ok {
		response.ValidationErrors = s.validationErrors(lang, errs)
	}

	return c.JSON(nestedResp.httpStatusCode, response)
}

func (s *Service) newResponse(lang msgLang, opt *parsedOpt) *Response {
	if opt.message == nil {
		opt.message = unknownMessage
	}

	txt := opt.message.Text(lang)
	if len(opt.args) != 0 {
		txt = fmt.Sprintf(txt, opt.args...)
	}

	return &Response{
		Code:             opt.message.Code,
		Message:          txt,
		AdditionalFields: opt.additionalFields,
	}
}

func (s *Service) TextFromMsgKey(key string) string {
	return s.msg(key, en)
}

func (s *Service) msg(key string, lang msgLang) string {
	if msg, ok := s.messages[key]; ok {
		msg.Text(lang)
	}

	return unknownMessage.Text(lang)
}
