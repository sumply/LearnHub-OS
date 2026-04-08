package create_quiz

import (
	"encoding/json"
	"server/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestQuestion_UnwrapJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonMap     map[string]any
		expected    func(*testing.T, domain.IQuestion, error)
		expectError bool
	}{
		{
			name: "single type",
			jsonMap: map[string]any{
				"text":  "question title",
				"score": 1,
				"type":  "single",
				"details": map[string]any{
					"correct": "opt1",
					"options": []string{"opt1", "opt2"},
				},
			},
			expected: func(t *testing.T, i domain.IQuestion, err error) {
				assert.IsType(t, &domain.SingleQuestion{}, i)

				single := i.(*domain.SingleQuestion)
				assert.Equal(t, "question title", single.Title)
				assert.Equal(t, 1, single.Score())
				assert.Equal(t, "opt1", single.Correct)
				assert.Equal(t, []string{"opt1", "opt2"}, single.Options)
			},
			expectError: false,
		},
		{
			name: "multiple type",
			jsonMap: map[string]any{
				"text":  "question title",
				"score": 1,
				"type":  "multiple",
				"details": map[string]any{
					"correct": []string{"opt1", "opt2"},
					"options": []string{"opt1", "opt2", "opt3"},
				},
			},
			expected: func(t *testing.T, i domain.IQuestion, err error) {
				assert.IsType(t, &domain.MultipleQuestion{}, i)

				single := i.(*domain.MultipleQuestion)
				assert.Equal(t, "question title", single.Title)
				assert.Equal(t, 1, single.Score())
				assert.Equal(t, []string{"opt1", "opt2"}, single.Correct)
				assert.Equal(t, []string{"opt1", "opt2", "opt3"}, single.Options)
			},
			expectError: false,
		},
		{
			name: "numeric type",
			jsonMap: map[string]any{
				"text":  "question title",
				"score": 1,
				"type":  "numeric",
				"details": map[string]any{
					"correct": float32(1),
				},
			},
			expected: func(t *testing.T, i domain.IQuestion, err error) {
				assert.IsType(t, &domain.NumericQuestion{}, i)

				single := i.(*domain.NumericQuestion)
				assert.Equal(t, "question title", single.Title)
				assert.Equal(t, 1, single.Score())
				assert.Equal(t, float32(1), single.Correct)
			},
			expectError: false,
		},
		{
			name: "invalid type",
			jsonMap: map[string]any{
				"text":  "question title",
				"score": 1,
				"type":  "invalid",
			},
			expectError: true,
		}, {
			name: "got domain error",
			jsonMap: map[string]any{
				"text":  "",
				"score": 1,
				"type":  "single",
				"details": map[string]any{
					"correct": "opt1",
					"options": []string{"opt1", "opt2"},
				},
			},
			expected: func(t *testing.T, i domain.IQuestion, err error) {
				assert.IsType(t, &domain.SingleQuestion{}, i)
				assert.Error(t, err)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.jsonMap)
			assert.NoError(t, err)

			req := &RequestQuestion{}
			err = json.Unmarshal(data, req)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.expected != nil {
					q, err := req.Unpack()
					tt.expected(t, q, err)
				}
			}
		})
	}
}
