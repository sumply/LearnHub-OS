package jsonb

import (
	"encoding/json"
	"server/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestQuestionRow_UnmarshalJSON(t *testing.T) {
	var preparedID = uuid.New()
	tests := []struct {
		name        string
		json        map[string]any
		result      func(*testing.T, domain.IQuestion, error)
		expectError bool
	}{
		{
			name: "single question",
			json: map[string]any{
				"id":             preparedID,
				"title":          "Question",
				"score":          3,
				"type":           "single",
				"single_correct": "yes",
				"single_options": []string{"no", "yes"},
			},
			result: func(t *testing.T, i domain.IQuestion, err error) {
				assert.IsType(t, &domain.SingleQuestion{}, i)
				assert.NoError(t, err)
				single := i.(*domain.SingleQuestion)

				assert.Equal(t, preparedID, single.ID)
				assert.Equal(t, "Question", single.Title)
				assert.Equal(t, 3, single.Score())
				assert.Equal(t, "yes", single.Correct)
				assert.Equal(t, []string{"no", "yes"}, single.Options)
			},
			expectError: false,
		},
		{
			name: "multiple question",
			json: map[string]any{
				"id":               preparedID,
				"title":            "Question",
				"score":            3,
				"type":             "multiple",
				"multiple_correct": []string{"yes", "maybe"},
				"multiple_options": []string{"no", "yes", "maybe"},
			},
			result: func(t *testing.T, i domain.IQuestion, err error) {
				assert.IsType(t, &domain.MultipleQuestion{}, i)
				assert.NoError(t, err)
				single := i.(*domain.MultipleQuestion)

				assert.Equal(t, preparedID, single.ID)
				assert.Equal(t, "Question", single.Title)
				assert.Equal(t, 3, single.Score())
				assert.Equal(t, []string{"yes", "maybe"}, single.Correct)
				assert.Equal(t, []string{"no", "yes", "maybe"}, single.Options)
			},
			expectError: false,
		},
		{
			name: "numeric question",
			json: map[string]any{
				"id":              preparedID,
				"title":           "Question",
				"score":           3,
				"type":            "numeric",
				"numeric_correct": 34.2,
			},
			result: func(t *testing.T, i domain.IQuestion, err error) {
				assert.IsType(t, &domain.NumericQuestion{}, i)
				assert.NoError(t, err)
				single := i.(*domain.NumericQuestion)

				assert.Equal(t, preparedID, single.ID)
				assert.Equal(t, "Question", single.Title)
				assert.Equal(t, 3, single.Score())
				assert.Equal(t, float32(34.2), single.Correct)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.json)
			assert.NoError(t, err)

			qr := &QuestionRow{}
			err = json.Unmarshal(data, qr)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.result != nil {
					d, e := qr.Unpack()
					tt.result(t, d, e)
				}
			}
		})
	}
}
