package httprpc

import (
	"errors"

	"github.com/olvrng/rbot/be/pkg/xerrors"
)

type Message interface{}

type ErrorInterface interface {
	// Code is of the valid error codes.
	Code() xerrors.Code

	// Msg returns a human-readable, unstructured messages describing the error.
	Msg() string

	// WithMeta returns a copy of the Error with the given key-value pair attached
	// as metadata. If the key is already set, it is overwritten.
	WithMeta(key string, val string) ErrorInterface

	// Meta returns the stored value for the given key. If the key has no set
	// value, Meta returns an empty string. There is no way to distinguish between
	// an unset value and an explicit empty string.
	Meta(key string) string

	// MetaMap returns the complete key-value metadata map stored on the error.
	MetaMap() map[string]string

	// Error returns a string of the form "twirp error <Type>: <Msg>"
	Error() string
}

type TwError interface {
	ErrorInterface
	Cause() error
	OrigFile() string
	OrigLine() int
}

type twError struct {
	err  *xerrors.APIError
	meta map[string]string
}

func (t twError) Code() xerrors.Code {
	return t.err.Code
}

func (t twError) Msg() string {
	return t.err.Message
}

func (t twError) Meta(key string) string {
	meta := t.meta
	if meta != nil {
		return meta[key]
	}
	return ""
}

func (t twError) WithMeta(key string, val string) ErrorInterface {
	t.meta[key] = val
	return t
}

func (t twError) MetaMap() map[string]string {
	return t.meta
}

func (t twError) Error() string {
	return t.err.Error()
}

func (t twError) Cause() error {
	if t.err.Err != nil {
		return t.err.Err
	}
	if t.err.Original != "" {
		return errors.New(t.err.Original)
	}
	return nil
}

func (t twError) OrigFile() string {
	return t.err.OrigFile
}

func (t twError) OrigLine() int {
	return t.err.OrigLine
}
