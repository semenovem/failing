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
				Translations: map[msgLang]string{en: "english_text"},
			},
		},
		logger: &testLogger{},
	}

	t.Run("1", func(t *testing.T) {
		fields := map[string]interface{}{
			"test-key": "test-value",
		}
		opt := serv.parseOpts([]interface{}{fields})

		assert.Equal(t, map[string]interface{}{
			"test-key": "test-value",
		}, opt.additionalFields)
	})

	t.Run("2", func(t *testing.T) {
		err := fmt.Errorf("test")
		opt := serv.parseOpts([]interface{}{err})

		assert.Equal(t, "test", opt.err.Error())
	})

	t.Run("3", func(t *testing.T) {
		args := Args{"one", "two"}
		opt := serv.parseOpts([]interface{}{args})

		assert.Equal(t, Args{"one", "two"}, opt.args)
	})

	t.Run("4", func(t *testing.T) {
		message := Message{
			Code:         5006,
			DefaultText:  "Дефолтный текст2",
			Translations: map[msgLang]string{en: "english2"},
		}

		opt := serv.parseOpts([]interface{}{message})

		assert.Equal(t, &Message{
			Code:         5006,
			DefaultText:  "Дефолтный текст2",
			Translations: map[msgLang]string{en: "english2"},
		}, opt.message)
	})

	t.Run("5", func(t *testing.T) {
		opt := serv.parseOpts([]interface{}{msgKey})

		assert.Equal(t, &Message{
			Code:         4002,
			DefaultText:  "Дефолтный текст",
			Translations: map[msgLang]string{en: "english_text"},
		}, opt.message)
	})
}
