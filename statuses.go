package failing

import "net/http"

var statuses = map[int]*Message{
	http.StatusBadRequest: {
		Code:        http.StatusBadRequest,
		DefaultText: "Неправильный запрос",
		Translations: map[Lang]string{
			EN: http.StatusText(http.StatusBadRequest),
		},
	},
	http.StatusUnauthorized: {
		Code:        http.StatusUnauthorized,
		DefaultText: "Вы не авторизованы",
		Translations: map[Lang]string{
			EN: http.StatusText(http.StatusUnauthorized),
		},
	},
	http.StatusForbidden: {
		Code:        http.StatusForbidden,
		DefaultText: "Запрещено",
		Translations: map[Lang]string{
			EN: http.StatusText(http.StatusForbidden),
		},
	},
	http.StatusNotFound: {
		Code:        http.StatusNotFound,
		DefaultText: "Объект не найден",
		Translations: map[Lang]string{
			EN: http.StatusText(http.StatusNotFound),
		},
	},
	http.StatusTooManyRequests: {
		Code:        http.StatusTooManyRequests,
		DefaultText: "Слишком много запросов",
		Translations: map[Lang]string{
			EN: http.StatusText(http.StatusTooManyRequests),
		},
	},
	http.StatusInternalServerError: {
		Code:        http.StatusInternalServerError,
		DefaultText: "Внутренняя ошибка сервера",
		Translations: map[Lang]string{
			EN: http.StatusText(http.StatusInternalServerError),
		},
	},
}
