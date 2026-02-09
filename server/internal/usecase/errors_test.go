package usecase

import (
	"errors"
	"server/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidationError(t *testing.T) {
	domainErr := getDomainError(t)
	validErr := NewValidationError(domainErr)

	assert.Equal(t, domainErr, validErr.Err())
	assert.Equal(t, "user validation", validErr.Event())

	anotherErr := errors.New("another")
	validErr = NewValidationError(anotherErr)
	assert.Equal(t, anotherErr, validErr.Err())
	assert.Equal(t, "validation", validErr.Event())
}

func getDomainError(t *testing.T) error {
	_, err := domain.NewUser("", "", "", "", "")

	assert.Error(t, err)

	var target *domain.Error
	assert.ErrorAs(t, err, &target)

	return err
}
