package errors

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

var (
	Wrapf = errors.Wrapf
)

const (
	DefaultInternalServerErrorID = "internal_server_error"
	DefaultMethodNotAllowedID    = "method_not_allowed"
	DefaultNotFoundID            = "not_found"
)

type Error struct {
	ID     string `json:"id,omitempty"`
	Code   int32  `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
	Status string `json:"status,omitempty"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var merr *Error
	if errors.As(err, &merr) {
		return merr, true
	}
	return nil, false
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if verr, ok := err.(*Error); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

func InternalServerError(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultInternalServerErrorID
	}
	return &Error{
		ID:     id,
		Code:   http.StatusInternalServerError,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusInternalServerError),
	}
}

func MethodNotAllowed(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultMethodNotAllowedID
	}
	return &Error{
		ID:     id,
		Code:   http.StatusMethodNotAllowed,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusMethodNotAllowed),
	}
}

func NotFound(id, format string, a ...interface{}) error {
	if id == "" {
		id = DefaultNotFoundID
	}
	return &Error{
		ID:     id,
		Code:   http.StatusNotFound,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusNotFound),
	}
}
