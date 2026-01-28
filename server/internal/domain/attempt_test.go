package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewAttempt(t *testing.T) {
	tests := []struct {
		name          string
		param         attemptParam
		expectedError error
	}{
		{
			name: "quiz is nil",
			param: attemptParam{
				quiz:   nil,
				userID: uuid.New(),
			},
			expectedError: ErrInvalid,
		},
		{
			name: "valid",
			param: attemptParam{
				quiz:   &Quiz{},
				userID: uuid.New(),
			},
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewAttempt(test.param.quiz, test.param.userID)
			if !errors.Is(err, test.expectedError) {
				t.Errorf("expected error: %s; got: %s", test.expectedError.Error(), err)
			}
		})
	}
}

func TestNewAttemptCreateAnswers(t *testing.T) {
	questions := []QuestionAggregate{
		{
			ID: uuid.New(),
		},
		{
			ID: uuid.New(),
		},
		{
			ID: uuid.New(),
		},
	}
	attempt, err := NewAttempt(&Quiz{
		Content: questions,
	}, uuid.New())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	for i, answer := range attempt.Answers {
		if questions[i].ID != answer.QuestionID {
			t.Errorf("incorrect answer; question id: %v; question id in answer: %v", questions[i].ID, answer.QuestionID)
		}
	}
}

type attemptParam struct {
	quiz   *Quiz
	userID uuid.UUID
}
