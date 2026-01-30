package dto

import (
	"encoding/json"
	"server/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalQuestion(t *testing.T) {
	request := []map[string]any{
		{
			"text":  "Question 1",
			"score": 1,
			"type":  "single",
			"details": map[string]any{
				"options": []string{"exception", "error", "throw", "fault"},
				"correct": "error",
			},
		},
		{
			"text":  "Question 2",
			"score": 1,
			"type":  "multiple",
			"details": map[string]any{
				"options": []string{"exception", "error", "throw", "fault"},
				"correct": []string{"error"},
			},
		},
		{
			"text":  "Question 3",
			"score": 1,
			"type":  "numeric",
			"details": map[string]any{
				"correct": 12.5,
			},
		},
	}
	data, err := json.Marshal(&request)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var input []Question
	if err := json.Unmarshal(data, &input); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if details, ok := input[0].Details.Domain.(*domain.SingleChoiceQuestion); !ok {
		t.Errorf("unexpected details type: %T", details)
	}

	if details, ok := input[1].Details.Domain.(*domain.MultipleChoiceQuestion); !ok {
		t.Errorf("unexpected details type: %T", details)
	}

	if details, ok := input[2].Details.Domain.(*domain.NumericQuestion); !ok {
		t.Errorf("unexpected details type: %T", details)
	}
}

func TestInvalidUnmarshalQuestion1(t *testing.T) {
	request := []map[string]any{
		{
			"text":  "Question 1",
			"score": 1,
			"type":  "invalid",
			"details": map[string]any{
				"options": []string{"exception", "error", "throw", "fault"},
				"correct": "error",
			},
		},
	}
	data, err := json.Marshal(&request)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var input []Question
	if err := json.Unmarshal(data, &input); err == nil {
		t.Error("expected error but got nil")
	}
}

func TestInvalidUnmarshalQuestion2(t *testing.T) {
	request := []map[string]any{
		{
			"text":  "Question 1",
			"score": 1,
			"type":  "single",
			"details": map[string]any{
				"options": []string{"exception", "error", "throw", "fault"},
				"correct": "invalid",
			},
		},
	}
	data, err := json.Marshal(&request)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var input []Question
	if err := json.Unmarshal(data, &input); err == nil {
		t.Error("expected error but got nil")
	}
}

func TestMarshalQuestion(t *testing.T) {
	singleMap := map[string]any{
		"text":  "Who's meowing?",
		"score": 1,
		"type":  domain.TypeSingleChoice,
		"details": map[string]any{
			"options": []string{"Cat", "Dog"},
			"correct": "Cat",
		},
	}
	multipleMap := map[string]any{
		"text":  "What programming languages exist?",
		"score": 2,
		"type":  domain.TypeMultipleChoice,
		"details": map[string]any{
			"options": []string{"C++", "Golang", "Rfgbv"},
			"correct": []string{"C++", "Golang"},
		},
	}
	numericMap := map[string]any{
		"text":  "2 + 2?",
		"score": 3,
		"type":  domain.TypeNumeric,
		"details": map[string]any{
			"correct": 4,
		},
	}
	tests := []struct {
		name               string
		dto                Question
		expectedMarshaling map[string]any
	}{
		{
			name: "single domain",
			dto: Question{
				Text:  "Who's meowing?",
				Score: 1,
				Details: QuestionDetails{
					Domain: &domain.SingleChoiceQuestion{
						Options: []string{"Cat", "Dog"},
						Correct: "Cat",
					},
				},
			},
			expectedMarshaling: singleMap,
		},
		{
			name: "multiple domain",
			dto: Question{
				Text:  "What programming languages exist?",
				Score: 2,
				Details: QuestionDetails{
					Domain: &domain.MultipleChoiceQuestion{
						Options: []string{"C++", "Golang", "Rfgbv"},
						Correct: []string{"C++", "Golang"},
					},
				},
			},
			expectedMarshaling: multipleMap,
		},
		{
			name: "numeric domain",
			dto: Question{
				Text:  "2 + 2?",
				Score: 3,
				Details: QuestionDetails{
					Domain: &domain.NumericQuestion{
						Correct: 4,
					},
				},
			},
			expectedMarshaling: numericMap,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := json.Marshal(&test.dto)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			expected, err := json.Marshal(&test.expectedMarshaling)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.JSONEqf(t, string(expected), string(actual), "invalid marshaling.\nEXPECTED:\n%v\nGOT:\n%v", expected, actual)
		})
	}
}

func TestInvalidMarshalQuestion(t *testing.T) {
	tests := []struct {
		name string
		dto  Question
	}{
		{
			name: "details is nil",
			dto: Question{
				Text:  "text",
				Score: 1,
				Details: QuestionDetails{
					Domain: nil,
				},
			},
		},
		{
			name: "unexpected details",
			dto: Question{
				Text:  "text",
				Score: 1,
				Details: QuestionDetails{
					Domain: &unexpectedDetails{},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := json.Marshal(&test.dto)
			if err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

type unexpectedDetails struct{}

func (i *unexpectedDetails) Variant() domain.QuestionType {
	return "fff"
}

func (i *unexpectedDetails) CheckAnswer(a any) (bool, error) {
	return true, nil
}

func (i *unexpectedDetails) Validate() error {
	return nil
}
