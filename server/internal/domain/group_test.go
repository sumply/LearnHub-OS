package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGroupTrim(t *testing.T) {
	group, err := NewGroup("   Name    ")
	assert.NoError(t, err)
	assert.Equal(t, group.Name, "Name")
}

func TestNewGroup(t *testing.T) {
	tests := []struct {
		name        string
		groupName   string
		expectError bool
	}{
		{
			name:        "name is empty",
			groupName:   "",
			expectError: true,
		},
		{
			name:        "valid",
			groupName:   "Name",
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewGroup(test.groupName)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
