package jsonb

import (
	"encoding/json"
	"server/internal/domain"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidQuestionDetailsMarshalJSON(t *testing.T) {
	single, err := domain.NewSingleChoiceQuestion(
		[]string{"opt1", "opt2"},
		"opt1",
	)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := json.Marshal(QuestionDetails{
		Domain: &single,
	})
	if err != nil {
		t.Fatal(err)
	}
	expectedMap := map[string]any{
		"details": map[string]any{
			"Options": single.Options,
			"Correct": single.Correct,
		},
		"variant": domain.TypeSingleChoice,
	}
	expected, err := json.Marshal(&expectedMap)
	if err != nil {
		t.Fatal(err)
	}
	assert.JSONEqf(t, string(expected), string(actual), "invalid marshaling.\nEXPECTED:\n%v\nGOT:\n%v", expected, actual)
}

func TestValidQuestionDetailsUnmarshalJSON(t *testing.T) {
	var (
		options = []string{"opt1", "opt2"}
		correct = "opt1"
	)
	const (
		variant = domain.TypeSingleChoice
	)
	actualMap := map[string]any{
		"details": map[string]any{
			"Options": options,
			"Correct": correct,
		},
		"variant": variant,
	}

	actual, err := json.Marshal(actualMap)
	if err != nil {
		t.Fatal(err)
	}

	var details QuestionDetails
	if err := json.Unmarshal(actual, &details); err != nil {
		t.Error(err)
	}

	if details.Domain == nil {
		t.Fatalf("domain is nil")
	}

	expected := &domain.SingleChoiceQuestion{
		Correct: correct,
		Options: options,
	}

	assert.Equalf(t, expected, details.Domain, "unexpected struct type")
}

func TestValidQuizQuestionAGGUnmarshalJSON(t *testing.T) {
	var (
		correct = "opt1"
		options = []string{correct, "opt2"}
		id      = uuid.New()
		quizID  = uuid.New()
	)
	const (
		variant = domain.TypeSingleChoice
		title   = "title"
		score   = 1
	)
	unpreparedJSON := map[string]any{
		"id":      id,
		"quiz_id": quizID,
		"title":   title,
		"score":   score,
		"details": map[string]any{
			"details": map[string]any{
				"Options": options,
				"Correct": correct,
			},
			"variant": variant,
		},
	}
	preparedJSON, err := json.MarshalIndent(&unpreparedJSON, "", "\t")
	if err != nil {
		t.Fatal(err)
	}

	var jsonAgg QuizQuestionAGG
	if err := json.Unmarshal(preparedJSON, &jsonAgg); err != nil {
		t.Fatal(err)
	}

	expectedJsonAgg := QuizQuestionAGG{
		ID:     id,
		QuizID: quizID,
		Title:  title,
		Details: QuestionDetails{
			Domain: &domain.SingleChoiceQuestion{
				Options: options,
				Correct: correct,
			},
		},
		Score: score,
	}

	assert.Equal(t, expectedJsonAgg, jsonAgg)
}
