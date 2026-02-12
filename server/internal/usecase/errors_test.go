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

	assert.ErrorAs(t, validErr, &domainErr)
	assert.Equal(t, "user validation", validErr.Error())

	anotherErr := errors.New("another")
	validErr = NewValidationError(anotherErr)

	assert.ErrorAs(t, validErr, &anotherErr)
	assert.Equal(t, "validation", validErr.Error())
}

func getDomainError(t *testing.T) *domain.Error {
	_, err := domain.NewUser("", "", "", "", "")

	assert.Error(t, err)

	var target *domain.Error
	assert.ErrorAs(t, err, &target)

	return target
}
