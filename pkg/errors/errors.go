package errors

import (
	"fmt"
	"github.com/pkg/errors"
)

type ErrorType int

const (
	NoType = ErrorType(iota)
	ValidationError
	ParsingError
	InvalidArgument
	InternalError
	RuntimeError
)

type customError struct {
	errorType     ErrorType
	originalError error
}

func (err *customError) Error() string {
	switch err.errorType {
	case ValidationError:
		return "validation error:" + err.originalError.Error()
	case ParsingError:
		return "parsing error:" + err.originalError.Error()
	case InvalidArgument:
		return "invalid argument error:" + err.originalError.Error()
	case RuntimeError:
		return "runtime error:" + err.originalError.Error()
	case InternalError:
		return "internal error:" + err.originalError.Error() + " (please report this error)"
	default:
		return err.originalError.Error()
	}
}

func (errType ErrorType) New(msg string) error {
	return &customError{
		errorType:     errType,
		originalError: errors.New(msg),
	}
}

func (errType ErrorType) Newf(format string, args ...interface{}) error {
	return &customError{
		errorType:     errType,
		originalError: fmt.Errorf(format, args...),
	}
}

func (errType ErrorType) Wrap(err error, msg string) error {
	return &customError{errorType: errType, originalError: errors.Wrapf(err, msg)}
}

func (errType ErrorType) Wrapf(err error, format string, args ...interface{}) error {
	return &customError{errorType: errType, originalError: errors.Wrapf(err, format, args...)}
}

type withMessage struct {
	message       string
	originalError error
}

func (err *withMessage) Error() string {
	return fmt.Sprintf("%s:\n\t%s", err.message, err.Error())
}

func WithMessage(err error, message string) error {
	return &withMessage{
		originalError: err,
		message:       message,
	}
}

func GetType(err error) ErrorType {
	if customErr, ok := err.(*customError); ok {
		return customErr.errorType
	}
	return NoType
}

func Unwrap(err error) error {
	if customErr, ok := err.(*customError); ok {
		return customErr.originalError
	}
	return err
}
