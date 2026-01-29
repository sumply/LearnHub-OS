package create_quiz

import (
	"encoding/json"
	"server/internal/domain"
	"testing"
)

func TestUnmarshalInputQuestion(t *testing.T) {
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

	var input []InputQuestion
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

func TestInvalidUnmarshalInputQuestion1(t *testing.T) {
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

	var input []InputQuestion
	if err := json.Unmarshal(data, &input); err == nil {
		t.Error("expected error but got nil")
	}
}

func TestInvalidUnmarshalInputQuestion2(t *testing.T) {
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

	var input []InputQuestion
	if err := json.Unmarshal(data, &input); err == nil {
		t.Error("expected error but got nil")
	}
}
