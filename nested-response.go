package failing

import (
	"fmt"
	"net/http"
)

type Nested interface {
	Message() string
	getOpts() []interface{}
	getValidationErr() error
	getHTTPStatusCode() int
}

type nestedResponse struct {
	httpStatusCode int
	opts           []interface{}
	validationErr  error
}

func NewNested(httpStatusCode int, opts ...interface{}) Nested {
	return &nestedResponse{
		httpStatusCode: httpStatusCode,
		opts:           opts,
	}
}

func NewNestedValidation(httpStatusCode int, err error, opts ...interface{}) Nested {
	return &nestedResponse{
		httpStatusCode: httpStatusCode,
		opts:           opts,
		validationErr:  err,
	}
}

func (n *nestedResponse) Message() string {
	if n.validationErr != nil {
		return invalidRequestMessage.Text(EN)
	}

	for _, v := range n.opts {
		if err, ok := v.(error); ok {
			return err.Error()
		}
	}

	if ret := http.StatusText(n.httpStatusCode); ret != "" {
		return ret
	}

	return fmt.Sprintf("%v", n)
}

func (n *nestedResponse) getOpts() []interface{} {
	return n.opts
}

func (n *nestedResponse) getValidationErr() error {
	return n.validationErr
}

func (n *nestedResponse) getHTTPStatusCode() int {
	return n.httpStatusCode
}
