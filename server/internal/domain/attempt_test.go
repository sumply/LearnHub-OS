package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewInvalidAttempt(t *testing.T) {
	tests := []struct {
		name  string
		param attemptParam
	}{
		{
			name: "quiz is nil",
			param: attemptParam{
				quiz:   nil,
				userID: uuid.New(),
			},
		},
		{
			name: "userID is empty",
			param: attemptParam{
				quiz:   &Quiz{},
				userID: uuid.Nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewAttempt(test.param.quiz, test.param.userID)
			assert.Error(t, err)
		})
	}
}

func TestAttemptAnswers(t *testing.T) {
	test := struct {
		name  string
		param attemptParam
	}{
		name: "answers length",
		param: attemptParam{
			quiz: &Quiz{
				Questions: []Question{
					{
						ID: uuid.New(),
					},
					{
						ID: uuid.New(),
					},
					{
						ID: uuid.New(),
					},
				},
			},
			userID: uuid.New(),
		},
	}
	t.Run(test.name, func(t *testing.T) {
		attempt, err := NewAttempt(test.param.quiz, uuid.New())
		assert.NoError(t, err)

		for i, answer := range attempt.Answers {
			if answer.AttemptID != attempt.ID {
				t.Error("answer attempt id is not attempt id")
			}
			if answer.QuestionID != test.param.quiz.Questions[i].ID {
				t.Error("question id is invalid")
			}
		}
	})
}

type attemptParam struct {
	quiz   *Quiz
	userID uuid.UUID
}
