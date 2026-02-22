package usecase

import (
	"errors"
	"server/internal/domain"
)

type ValidationError struct {
	event string
	err   error
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
	return e.event
}

func (e *ValidationError) Unwrap() error {
	return e.err
}

type AuthError struct {
	event string
	err   error
}

func NewAuthError(err error) *AuthError {
	return &AuthError{
		event: "permission denied",
		err:   err,
	}
}

func (a *AuthError) Error() string {
	return a.event
}

func (a *AuthError) Unwrap() error {
	return a.err
}
