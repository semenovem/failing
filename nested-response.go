package failing

import (
	"fmt"
	"net/http"
)

type NestedResponse struct {
	httpStatusCode int
	opts           []interface{}
}

func NewError(httpStatusCode int, opts ...interface{}) *NestedResponse {
	return &NestedResponse{
		httpStatusCode: httpStatusCode,
		opts:           opts,
	}
}

func (n *NestedResponse) Message(service *Service) string {
	for _, v := range n.opts {
		if err, ok := v.(error); ok {
			return err.Error()
		}
	}
	// TODO вынести в метод - конвертер
	for _, v := range n.opts {
		if msgKey, ok := v.(string); ok {
			if msg, ok := service.messages[msgKey]; ok {
				return msg.DefaultText
			}
		}
	}

	if ret := http.StatusText(n.httpStatusCode); ret != "" {
		return ret
	}

	return fmt.Sprintf("%v", n)
}
