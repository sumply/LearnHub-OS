package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorString(t *testing.T) {
	err := error(newError(map[string]error{
		"name": ErrInvalid,
		"age":  ErrInvalid,
	}))

	expected := fmt.Sprintf("name: %s; age: %s", ErrInvalid, ErrInvalid)

	assert.Equal(t, expected, err.Error())
}
