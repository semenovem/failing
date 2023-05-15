package failing

type Message struct {
	Code         int             // Код ошибки
	DefaultText  string          // Текст по умолчанию
	Translations map[Lang]string // Переводы
}

func (m *Message) Text(language Lang) string {
	if d, ok := m.Translations[language]; ok {
		return d
	}

	return m.DefaultText
}
