package failing

const (
	headerLanguageName = "Accept-Language"
	fieldNameErr       = "__error__" // Имя поля в которое добавляется содержимое ошибки в dev режиме

	RU Lang = "ru-RU"
	EN Lang = "en-EN"

	unknownText = "unknownText"
)

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

type Response struct {
	Code             int                    `json:"code"`
	Message          string                 `json:"message"`
	ValidationErrors []*ValidationError     `json:"validation_errors,omitempty"`
	AdditionalFields map[string]interface{} `json:"additional_fields,omitempty"`
}

type Lang string

var (
	unknownMessage = &Message{
		Code:        1000,
		DefaultText: "Неизвестная ошибка",
		Translations: map[Lang]string{
			EN: "Unknown error",
		},
	}

	invalidRequestMessage = &Message{
		Code:        1001,
		DefaultText: "Невалидные параметры запроса",
		Translations: map[Lang]string{
			EN: "Invalid request parameters",
		},
	}
)

// Args аргументы для шаблона fmt.Sprint(...)
type Args []interface{}

type logger interface {
	Errorf(template string, args ...interface{})
}
