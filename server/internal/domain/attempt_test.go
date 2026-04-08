package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnswerMaker_VisitSingleQuestion(t *testing.T) {
	question := &SingleQuestion{}
	maker := &answerMaker{}
	err := question.Accept(maker)
	assert.NoError(t, err)
	assert.IsType(t, &SingleAnswer{}, maker.maked)
}

func TestAnswerMaker_VisitMultipleQuestion(t *testing.T) {
	question := &MultipleQuestion{}
	maker := &answerMaker{}
	err := question.Accept(maker)
	assert.NoError(t, err)
	assert.IsType(t, &MultipleAnswer{}, maker.maked)
}

func TestAnswerMaker_VisitNumericQuestion(t *testing.T) {
	question := &NumericQuestion{}
	maker := &answerMaker{}
	err := question.Accept(maker)
	assert.NoError(t, err)
	assert.IsType(t, &NumericAnswer{}, maker.maked)
}

func TestSingleAnswer_Check(t *testing.T) {
	question, err := NewSingleQuestion("title", "yes", []string{"no", "yes"}, 1)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		answer      *SingleAnswer
		expectError bool
		isValid     bool
	}{
		{
			name: "answer to the wrong question",
			answer: func() *SingleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*SingleAnswer)
				answer.questionID = uuid.New()
				return answer
			}(),
			expectError: true,
			isValid:     false,
		},
		{
			name: "correct answer",
			answer: func() *SingleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*SingleAnswer)
				answer.Selected = "yes"
				return answer
			}(),
			expectError: false,
			isValid:     true,
		},
		{
			name: "incorrect answer",
			answer: func() *SingleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*SingleAnswer)
				answer.Selected = "no"
				return answer
			}(),
			expectError: false,
			isValid:     false,
		},
		{
			name: "selected answer is invalid",
			answer: func() *SingleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*SingleAnswer)
				answer.Selected = "fsdfaas"
				return answer
			}(),
			expectError: true,
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.answer.Check(question)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.isValid, tt.answer.isCorrect)
			}
		})
	}
}

func TestMultipleAnswer_Check(t *testing.T) {
	question, err := NewMultipleQuestion("title", []string{"yes", "maybe"}, []string{"no", "yes", "maybe"}, 1)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		answer      *MultipleAnswer
		expectError bool
		isValid     bool
	}{
		{
			name: "answer to the wrong question",
			answer: func() *MultipleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*MultipleAnswer)
				answer.questionID = uuid.New()
				return answer
			}(),
			expectError: true,
			isValid:     false,
		},
		{
			name: "correct answer",
			answer: func() *MultipleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*MultipleAnswer)
				answer.Selected = []string{"yes", "maybe"}
				return answer
			}(),
			expectError: false,
			isValid:     true,
		},
		{
			name: "incorrect answer",
			answer: func() *MultipleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*MultipleAnswer)
				answer.Selected = []string{"no"}
				return answer
			}(),
			expectError: false,
			isValid:     false,
		},
		{
			name: "half correct answer",
			answer: func() *MultipleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*MultipleAnswer)
				answer.Selected = []string{"yes"}
				return answer
			}(),
			expectError: false,
			isValid:     false,
		},
		{
			name: "selected answer is invalid",
			answer: func() *MultipleAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*MultipleAnswer)
				answer.Selected = []string{"mmmmmmm"}
				return answer
			}(),
			expectError: true,
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.answer.Check(question)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.isValid, tt.answer.isCorrect)
			}
		})
	}
}

func TestNumericAnswer_Check(t *testing.T) {
	question, err := NewNumericQuestion("title", 342.2, 1)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		answer      *NumericAnswer
		expectError bool
		isValid     bool
	}{
		{
			name: "answer to the wrong question",
			answer: func() *NumericAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*NumericAnswer)
				answer.questionID = uuid.New()
				return answer
			}(),
			expectError: true,
			isValid:     false,
		},
		{
			name: "correct answer",
			answer: func() *NumericAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*NumericAnswer)
				answer.Selected = 342.2
				return answer
			}(),
			expectError: false,
			isValid:     true,
		},
		{
			name: "incorrect answer",
			answer: func() *NumericAnswer {
				maker := new(answerMaker)
				err := question.Accept(maker)
				assert.NoError(t, err)
				answer := maker.maked.(*NumericAnswer)
				answer.Selected = 23
				return answer
			}(),
			expectError: false,
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.answer.Check(question)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tt.isValid, tt.answer.isCorrect)
			}
		})
	}
}

func TestNewAttempt(t *testing.T) {
	questions := []IQuestion{
		func() IQuestion {
			single, err := NewSingleQuestion("single", "yes", []string{"yes", "no"}, 1)
			assert.NoError(t, err)
			return single
		}(),
		func() IQuestion {
			mutliple, err := NewMultipleQuestion("multiple", []string{"yes", "maybe"}, []string{"yes", "no", "maybe"}, 1)
			assert.NoError(t, err)
			return mutliple
		}(),
		func() IQuestion {
			mutliple, err := NewNumericQuestion("numeric", 234.23, 1)
			assert.NoError(t, err)
			return mutliple
		}(),
	}
	quiz, err := NewQuiz(uuid.New(), uuid.New(), "Test", "Test", questions, uuid.UUIDs{uuid.New(), uuid.New()}, 3, nil)
	assert.NoError(t, err)
	questionMap := make(map[uuid.UUID]IQuestion)
	for _, q := range quiz.Questions {
		questionMap[q.ID()] = q
	}
	assert.Equal(t, len(questionMap), len(quiz.Questions))
	attempt, err := NewAttempt(quiz, uuid.New())
	assert.NoError(t, err)
	assert.Equal(t, len(questionMap), len(attempt.SelectedAnswers))
}
