package domain

import (
	"errors"
	"strings"
)

var (
	ErrInvalid  = errors.New("invaid")
	ErrValidate = errors.New("validate")
)

type Error struct {
	Fields map[string]error
}

func newError(fields map[string]error) *Error {
	return &Error{
		Fields: fields,
	}
}

func (e *Error) Error() string {
	var s string
	for field, err := range e.Fields {
		s += field + ": " + err.Error() + "; "
	}
	s = strings.TrimSuffix(s, "; ")
	return s
}
