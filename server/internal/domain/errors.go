package domain

import (
	"errors"
	"strings"
)

type Error struct {
	domain string
	data   []ErrorData
}

type ErrorData struct {
	Field string
	Error error
}

func NewError(domain string) *Error {
	return &Error{
		domain: domain,
	}
}

func (f *Error) Add(field string, err error) {
	f.data = append(f.data, ErrorData{
		Field: field,
		Error: err,
	})
}

func (f *Error) Empty() bool {
	return len(f.data) == 0
}

func (f *Error) Data() []ErrorData {
	return f.data
}

func (f *Error) Domain() string {
	return f.domain
}

func (e *Error) Error() string {
	var sb strings.Builder
	sb.WriteString(e.domain)
	sb.WriteString("[")
	isAdded := false
	for i := range e.data {
		if isAdded {
			sb.WriteString("; ")
		} else {
			isAdded = true
		}
		sb.WriteString(e.data[i].Field)
		sb.WriteString(": ")
		sb.WriteString(e.data[i].Error.Error())
	}
	sb.WriteString("]")
	return sb.String()
}

func (e *Error) ToMap() map[string]any {
	fieldMap := make(map[string]any)
	for i := range e.data {
		var target *Error
		if errors.As(e.data[i].Error, &target) {
			fieldMap[e.data[i].Field] = target.ToMap()
		}
		fieldMap[e.data[i].Field] = e.data[i].Error.Error()
	}

	return map[string]any{
		e.domain: fieldMap,
	}
}
