package failing

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testLogger struct{}

func (l *testLogger) Errorf(template string, args ...interface{}) {
	fmt.Printf(template, args...)
}

func Test_parseOpts(t *testing.T) {
	t.Parallel()

	msgKey := "key_001122"

	serv := &Service{
		messages: map[string]*Message{
			msgKey: {
				Code:         4002,
				DefaultText:  "Дефолтный текст",
				Translations: map[Lang]string{EN: "english_text"},
			},
		},
		logger: &testLogger{},
	}

	servDev := &Service{
		messages: map[string]*Message{
			msgKey: {
				Code:         4002,
				DefaultText:  "Дефолтный текст",
				Translations: map[Lang]string{EN: "english_text"},
			},
		},
		logger: &testLogger{},
		isDev:  true,
	}

	t.Run("additionalFields", func(t *testing.T) {
		fields := map[string]interface{}{
			"test-key": "test-value",
		}

		assert.Equal(t, map[string]interface{}{
			"test-key": "test-value",
		}, serv.parseOpts([]interface{}{fields}).additionalFields)

		assert.Equal(t, map[string]interface{}{
			"test-key": "test-value",
		}, servDev.parseOpts([]interface{}{fields}).additionalFields)
	})

	t.Run("error", func(t *testing.T) {
		var (
			err    = fmt.Errorf("test321")
			fields = map[string]interface{}{
				"test-key": "test-value",
			}
		)

		assert.Equal(t, map[string]interface{}{
			"test-key": "test-value",
		}, serv.parseOpts([]interface{}{err, fields}).additionalFields)

		assert.Equal(t, map[string]interface{}{
			"test-key":  "test-value",
			"__error__": "test321",
		}, servDev.parseOpts([]interface{}{err, fields}).additionalFields)
	})

	t.Run("Args", func(t *testing.T) {
		args := Args{"one", "two"}

		assert.Equal(t, Args{"one", "two"}, serv.parseOpts([]interface{}{args}).args)
		assert.Equal(t, Args{"one", "two"}, servDev.parseOpts([]interface{}{args}).args)
	})

	t.Run("Message", func(t *testing.T) {
		message := Message{
			Code:         5006,
			DefaultText:  "Дефолтный текст2",
			Translations: map[Lang]string{EN: "english2"},
		}

		assert.Equal(t, &Message{
			Code:         5006,
			DefaultText:  "Дефолтный текст2",
			Translations: map[Lang]string{EN: "english2"},
		}, serv.parseOpts([]interface{}{message}).message)

		assert.Equal(t, &Message{
			Code:         5006,
			DefaultText:  "Дефолтный текст2",
			Translations: map[Lang]string{EN: "english2"},
		}, servDev.parseOpts([]interface{}{message}).message)
	})

	t.Run("message-key", func(t *testing.T) {
		assert.Equal(t, &Message{
			Code:         4002,
			DefaultText:  "Дефолтный текст",
			Translations: map[Lang]string{EN: "english_text"},
		}, serv.parseOpts([]interface{}{msgKey}).message)

		assert.Equal(t, &Message{
			Code:         4002,
			DefaultText:  "Дефолтный текст",
			Translations: map[Lang]string{EN: "english_text"},
		}, servDev.parseOpts([]interface{}{msgKey}).message)
	})
}
