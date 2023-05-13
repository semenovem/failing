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
			Translations: map[msgLang]string{
				en: "text in english",
			},
		},
	}

	validationMessageMap := map[string]string{
		"": "",
	}

	return New(&Config{
		IsDevMode:            false,
		TranslatorDefault:    nil,
		Translators:          nil,
		Messages:             messages,
		ValidationMessageMap: validationMessageMap,
		Logger:               nil,
	})
}

func genEchoRequest(lang msgLang) (echo.Context, *httptest.ResponseRecorder) {
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
		serv = getService()
		c, _ = genEchoRequest(ru)
	)
	//
	//t.Run("1", func(t *testing.T) {
	//	response := serv.NewResponse(c)
	//
	//	assert.Equal(t, 1000, response.Code)
	//	assert.Equal(t, unknownMessage.Text(ru), response.Message)
	//	assert.Nil(t, response.AdditionalFields)
	//	assert.Nil(t, response.AdditionalFields)
	//})

	t.Run("2", func(t *testing.T) {
		err := fmt.Errorf("test")

		response := serv.NewResponse(c, err)

		assert.Equal(t, 1000, response.Code)
		assert.Equal(t, unknownMessage.Text(ru), response.Message)
		assert.Nil(t, response.AdditionalFields)
		assert.Nil(t, response.AdditionalFields)
	})

}
