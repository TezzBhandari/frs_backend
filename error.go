package frs

import (
	"errors"
	"fmt"
)

const (
	EBADREQUEST   = "bad_request"
	EINVALID      = "invalid"
	EUNAUTHORIZED = "unauthorized"
	ENOTFOUND     = "not_found"
	EINTERNAL     = "internal_error"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("frs error: code: %s message: %s", e.Code, e.Message)
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	}

	if errors.As(err, &e) {
		return e.Code
	}

	return EINTERNAL
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	}

	if errors.As(err, &e) {
		return e.Message
	}

	return "internal error"
}

func Errorf(code string, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
