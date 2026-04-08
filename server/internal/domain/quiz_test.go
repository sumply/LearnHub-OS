package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type single struct {
	title   string
	correct string
	options []string
	score   int
}

func TestNewSingleQuestion(t *testing.T) {
	tests := []struct {
		name        string
		param       single
		expected    single
		expectError bool
	}{
		{
			name: "valid",
			param: single{
				title:   "title",
				correct: "opt1",
				options: []string{"opt1", "opt2"},
				score:   1,
			},
			expected: single{
				title:   "title",
				correct: "opt1",
				options: []string{"opt1", "opt2"},
				score:   1,
			},
			expectError: false,
		},
		{
			name: "trimmed fields",
			param: single{
				title:   "title    ",
				correct: "     opt1",
				options: []string{" opt1 ", "   opt2   "},
				score:   1,
			},
			expected: single{
				title:   "title",
				correct: "opt1",
				options: []string{"opt1", "opt2"},
				score:   1,
			},
			expectError: false,
		},
		{
			name: "title is empty",
			param: single{
				title:   "",
				correct: "     correct",
				options: []string{" opt1 ", "   opt2   "},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "correct is not contains in options",
			param: single{
				title:   "title",
				correct: "wrong",
				options: []string{"opt1", "opt2"},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "option is empty",
			param: single{
				title:   "title",
				correct: "opt1",
				options: []string{"opt1", ""},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "option is not unique",
			param: single{
				title:   "title",
				correct: "opt1",
				options: []string{"opt1", "opt1"},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "options is empty",
			param: single{
				title:   "title",
				correct: "opt1",
				options: nil,
				score:   1,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewSingleQuestion(
				tt.param.title,
				tt.param.correct,
				tt.param.options,
				tt.param.score,
			)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.title, q.Title)
				assert.Equal(t, tt.expected.correct, q.Correct)
				assert.Equal(t, tt.expected.options, q.Options)
				assert.Equal(t, tt.expected.score, q.score)
			}
		})
	}
}

type multiple struct {
	title   string
	correct []string
	options []string
	score   int
}

func TestMultipleQuestion(t *testing.T) {
	tests := []struct {
		name        string
		param       multiple
		expected    multiple
		expectError bool
	}{
		{
			name: "title is empty",
			param: multiple{
				title:   "",
				correct: []string{"opt1"},
				options: []string{"opt1"},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "correct is empty",
			param: multiple{
				title:   "title",
				correct: nil,
				options: []string{"opt1"},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "correct is not unique",
			param: multiple{
				title:   "title",
				correct: []string{"opt1", "opt1"},
				options: []string{"opt1", "opt2"},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "options is not unique",
			param: multiple{
				title:   "title",
				correct: []string{"opt1"},
				options: []string{"opt1", "opt1"},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "options is empty",
			param: multiple{
				title:   "title",
				correct: []string{"opt1"},
				options: nil,
				score:   1,
			},
			expectError: true,
		},
		{
			name: "correct is not contains in options",
			param: multiple{
				title:   "title",
				correct: []string{"opt1", "wrong"},
				options: []string{"opt1", "opt2"},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "option is empty",
			param: multiple{
				title:   "title",
				correct: []string{"opt1"},
				options: []string{"opt1", ""},
				score:   1,
			},
			expectError: true,
		},
		{
			name: "multiple correct",
			param: multiple{
				title:   "title",
				correct: []string{"opt1", "opt2"},
				options: []string{"opt1", "opt2", "opt3"},
				score:   1,
			},
			expected: multiple{
				title:   "title",
				correct: []string{"opt1", "opt2"},
				options: []string{"opt1", "opt2", "opt3"},
				score:   1,
			},
			expectError: false,
		},
		{
			name: "trimmed fields",
			param: multiple{
				title:   "     title ",
				correct: []string{" opt1   ", " opt2 "},
				options: []string{"   opt1  ", "   opt2   ", "  opt3   "},
				score:   1,
			},
			expected: multiple{
				title:   "title",
				correct: []string{"opt1", "opt2"},
				options: []string{"opt1", "opt2", "opt3"},
				score:   1,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewMultipleQuestion(
				tt.param.title,
				tt.param.correct,
				tt.param.options,
				tt.param.score,
			)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.title, q.Title)
				assert.Equal(t, tt.expected.correct, q.Correct)
				assert.Equal(t, tt.expected.options, q.Options)
				assert.Equal(t, tt.expected.score, q.score)
			}
		})
	}
}

type numeric struct {
	title   string
	correct float32
	score   int
}

func TestNumericQuestion(t *testing.T) {
	tests := []struct {
		name        string
		param       numeric
		expected    numeric
		expectError bool
	}{
		{
			name: "title is empty",
			param: numeric{
				title:   "",
				correct: 1,
				score:   1,
			},
			expectError: true,
		},
		{
			name: "valid",
			param: numeric{
				title:   "title",
				correct: 1,
				score:   1,
			},
			expected: numeric{
				title:   "title",
				correct: 1,
				score:   1,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewNumericQuestion(
				tt.param.title,
				tt.param.correct,
				tt.param.score,
			)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.title, q.Title)
				assert.Equal(t, tt.expected.correct, q.Correct)
				assert.Equal(t, tt.expected.score, q.score)
			}
		})
	}
}
