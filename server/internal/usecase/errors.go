package usecase

import (
	"errors"
	"fmt"
	"server/internal/domain"
)

type ValidationError struct {
	event string
	err   error
}

func (e *ValidationError) Event() string {
	return e.event
}

func (e *ValidationError) Err() error {
	return e.err
}

func NewValidationError(err error) *ValidationError {
	validErr := new(ValidationError)

	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		validErr.event = domainErr.Domain() + " validation"
	} else {
		validErr.event = "validation"
	}

	validErr.err = err

	return validErr
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("usecase(%s): %s", e.event, e.err.Error())
}
