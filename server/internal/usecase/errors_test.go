package usecase

import (
	"errors"
	"fmt"
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

func TestNewAuthError(t *testing.T) {
	detail := fmt.Errorf("user is not admin")
	auth := NewAuthError(detail)

	assert.Equal(t, "permission denied", auth.Error())
	assert.EqualError(t, auth.Unwrap(), detail.Error())
}
