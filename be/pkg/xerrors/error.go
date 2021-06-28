package xerrors

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// IError defines error interface returned by xerrors package
type IError interface {
	error
	IStack

	Format(st fmt.State, verb rune)
}

// IStack declares that the error has stacktrace. It's compatible with github.com/pkg/errors package.
type IStack interface {
	StackTrace() errors.StackTrace
}

type APIError struct {
	Code     Code
	Err      error
	Message  string
	Original string
	OrigFile string
	OrigLine int
	Stack    errors.StackTrace
	Trace    bool
	Trivial  bool
	Meta     map[string]string
}

func (e *APIError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	if e.Err != nil {
		b.WriteString(" cause=")
		b.WriteString(e.Err.Error())
	}
	if e.Original != "" {
		b.WriteString(" original=")
		b.WriteString(e.Original)
	}
	for k, v := range e.Meta {
		b.WriteByte(' ')
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
	}
	return b.String()
}

func (e *APIError) WithMeta(key string, value string) *APIError {
	if e.Meta == nil {
		e.Meta = make(map[string]string)
	}
	e.Meta[key] = value
	return e
}

func newError(trace bool, stack bool, code Code, message string, err error) *APIError {
	if !IsValidStandardErrorCode(code) {
		code = Unknown
	}
	if message == "" {
		message = DefaultErrorMessage(code)
	}

	if err != nil {
		// overwrite *Error
		if xerr, ok := err.(*APIError); ok {
			// keep original message
			if xerr.Original == "" {
				xerr.Original = xerr.Message
			}
			xerr.Code = code
			xerr.Message = message
			xerr.Trace = xerr.Trace || trace
			return xerr
		}
	}

	// always include the original location
	_, file, line, _ := runtime.Caller(2)
	xerr := &APIError{
		Err:      err,
		Code:     code,
		Message:  message,
		Original: "",
		OrigFile: file,
		OrigLine: line,
		Trace:    trace,
	}

	// wrap error with stacktrace
	if stack {
		xerr.Stack = errors.New("").(IStack).StackTrace()[2:]
	}
	return xerr
}

func Errorf(code Code, err error, message string, args ...interface{}) *APIError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return newError(false, true, code, message, err)
}

func GetCode(err error) Code {
	if err == nil {
		return NoError
	}
	xerr, ok := err.(*APIError)
	if !ok {
		return Unknown
	}
	return xerr.Code
}
