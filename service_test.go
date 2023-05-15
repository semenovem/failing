package failing

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getService() *Service {
	msgKeyErr1 := "1100 Текст на русском"

	messages := map[string]*Message{
		msgKeyErr1: {
			Code:        1100,
			DefaultText: "Текст на русском",
			Translations: map[Lang]string{
				EN: "text in english",
			},
		},
	}

	validationMessageMap := map[string]string{
		"": "",
	}

	return New(&Config{
		IsDevMode:            false,
		TranslatorDefault:    nil,
		Messages:             messages,
		ValidationMessageMap: validationMessageMap,
		Logger:               nil,
	})
}

func genEchoRequest(lang Lang) (echo.Context, *httptest.ResponseRecorder) {
	ho := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(headerLanguageName, string(lang))

	rec := httptest.NewRecorder()

	return ho.NewContext(req, rec), rec
}

func TestServiceNewResponse(t *testing.T) {
	t.Parallel()

	var (
		serv    = getService()
		servDev = getService()
	)

	servDev.isDev = true

	t.Run("1-unknown", func(t *testing.T) {
		var (
			err  = fmt.Errorf("test")
			c, _ = genEchoRequest(RU)
		)

		response := serv.NewResponse(c, err)

		assert.Equal(t, 1000, response.Code)
		assert.Equal(t, unknownMessage.Text(RU), response.Message)
		assert.Nil(t, response.AdditionalFields)
		assert.Nil(t, response.ValidationErrors)
	})

	t.Run("2-ru", func(t *testing.T) {
		var (
			err  = fmt.Errorf("test")
			c, _ = genEchoRequest(RU)
		)

		response := serv.NewResponse(c, err, "1100 Текст на русском")

		assert.Equal(t, 1100, response.Code)
		assert.Equal(t, "Текст на русском", response.Message)
		assert.Nil(t, response.AdditionalFields)
		assert.Nil(t, response.ValidationErrors)
	})

	t.Run("2-en", func(t *testing.T) {
		var (
			err  = fmt.Errorf("test")
			c, _ = genEchoRequest(EN)
		)

		response := serv.NewResponse(c, err, "1100 Текст на русском")

		assert.Equal(t, 1100, response.Code)
		assert.Equal(t, "text in english", response.Message)
		assert.Nil(t, response.AdditionalFields)
		assert.Nil(t, response.ValidationErrors)
	})

	t.Run("2-en-dev", func(t *testing.T) {
		var (
			err  = fmt.Errorf("test")
			c, _ = genEchoRequest(EN)
		)

		response := servDev.NewResponse(c, err, "1100 Текст на русском")

		assert.Equal(t, 1100, response.Code)
		assert.Equal(t, "text in english", response.Message)
		assert.Equal(t, map[string]interface{}{
			"__error__": "test",
		}, response.AdditionalFields)
		assert.Nil(t, response.ValidationErrors)
	})
}

func TestService_newResponse(t *testing.T) {
	t.Parallel()

	var serv = getService()

	t.Run("1", func(t *testing.T) {
		response := serv.newResponse(RU, &parsedOpt{})

		assert.Equal(t, 1000, response.Code)
		assert.Equal(t, unknownMessage.Text(RU), response.Message)
		assert.Nil(t, response.AdditionalFields)
		assert.Nil(t, response.AdditionalFields)
	})

	t.Run("2-ru", func(t *testing.T) {
		var (
			msg = Message{
				Code:        3545,
				DefaultText: "фыфывафывафываыва %v",
				Translations: map[Lang]string{
					EN: "asdfasfdasdfasdfas %v",
				},
			}

			opt = parsedOpt{
				message: &msg,
				args:    Args{"rrrttt"},
				additionalFields: map[string]interface{}{
					"field1": "value1",
				},
			}

			response = serv.newResponse(RU, &opt)
		)

		assert.Equal(t, 3545, response.Code)
		assert.Equal(t, "фыфывафывафываыва rrrttt", response.Message)
		assert.Equal(t, map[string]interface{}{
			"field1": "value1",
		}, response.AdditionalFields)
		assert.Nil(t, response.ValidationErrors)
	})

	t.Run("3-en", func(t *testing.T) {
		var (
			msg = Message{
				Code:        3545,
				DefaultText: "фыфывафывафываыва %v",
				Translations: map[Lang]string{
					EN: "asdfasfdasdfasdfas %v",
				},
			}

			opt = parsedOpt{
				message: &msg,
				args:    Args{"rrrttt"},
				additionalFields: map[string]interface{}{
					"field1": "value1",
				},
			}

			response = serv.newResponse(EN, &opt)
		)

		assert.Equal(t, 3545, response.Code)
		assert.Equal(t, "asdfasfdasdfasdfas rrrttt", response.Message)
		assert.Equal(t, map[string]interface{}{
			"field1": "value1",
		}, response.AdditionalFields)
		assert.Nil(t, response.ValidationErrors)
	})
}
