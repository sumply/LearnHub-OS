package get_quiz

import (
	"server/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseQuestion(t *testing.T) {
	var preparedID = uuid.New()
	tests := []struct {
		name        string
		param       domain.IQuestion
		expected    map[string]any
		expectError bool
	}{
		{
			name: "single question",
			param: func() *domain.SingleQuestion {
				s, err := domain.RestoreSingleQuestion(
					preparedID,
					"Question",
					"yes",
					[]string{"no", "yes"},
					1,
				)
				assert.NoError(t, err)
				return s
			}(),
			expected: map[string]any{
				"id":      preparedID,
				"title":   "Question",
				"score":   1,
				"correct": "yes",
				"options": []string{"no", "yes"},
				"type":    "single",
			},
			expectError: false,
		},
		{
			name: "multiple question",
			param: func() *domain.MultipleQuestion {
				s, err := domain.RestoreMultipleQuestion(
					preparedID,
					"Question",
					[]string{"yes", "maybe"},
					[]string{"no", "yes", "maybe"},
					1,
				)
				assert.NoError(t, err)
				return s
			}(),
			expected: map[string]any{
				"id":      preparedID,
				"title":   "Question",
				"score":   1,
				"correct": []string{"yes", "maybe"},
				"options": []string{"no", "yes", "maybe"},
				"type":    "multiple",
			},
			expectError: false,
		},
		{
			name: "numeric question",
			param: func() *domain.NumericQuestion {
				s, err := domain.RestoreNumericQuestion(
					preparedID,
					"Question",
					432.54,
					5,
				)
				assert.NoError(t, err)
				return s
			}(),
			expected: map[string]any{
				"id":      preparedID,
				"title":   "Question",
				"score":   5,
				"correct": float32(432.54),
				"type":    "numeric",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := NewResponseQuestion(tt.param)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp.data)
			}
		})
	}
}
