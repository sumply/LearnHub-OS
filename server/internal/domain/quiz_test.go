package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSingleChoiceQuestion(t *testing.T) {
	tests := []struct {
		name        string
		param       singleParam
		expectError bool
	}{
		{
			name: "options is empty",
			param: singleParam{
				options: nil,
				correct: "correct",
			},
			expectError: true,
		},
		{
			name: "correct is empty",
			param: singleParam{
				options: []string{"opt"},
				correct: "",
			},
			expectError: true,
		},
		{
			name: "answer is incorrect",
			param: singleParam{
				options: []string{"opt1", "opt2"},
				correct: "opt3",
			},
			expectError: true,
		},
		{
			name: "answer is correct",
			param: singleParam{
				options: []string{"opt1", "opt2"},
				correct: "opt1",
			},
			expectError: false,
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			_, err := NewSingleChoiceQuestion(tests[i].param.options, tests[i].param.correct)
			if !tests[i].expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if tests[i].expectError && err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

type singleParam struct {
	options []string
	correct string
}

func TestNewMultipleChoiceQuestion(t *testing.T) {
	tests := []struct {
		name        string
		param       multipleParam
		expectError bool
	}{
		{
			name: "options is empty",
			param: multipleParam{
				options: nil,
				correct: []string{"correct"},
			},
			expectError: true,
		},
		{
			name: "correct is empty",
			param: multipleParam{
				options: []string{"opt"},
				correct: nil,
			},
			expectError: true,
		},
		{
			name: "answer is incorrect",
			param: multipleParam{
				options: []string{"opt1", "opt2"},
				correct: []string{"opt3"},
			},
			expectError: true,
		},
		{
			name: "answer is correct",
			param: multipleParam{
				options: []string{"opt1", "opt2"},
				correct: []string{"opt1"},
			},
			expectError: false,
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			_, err := NewMultipleChoiceQuestion(tests[i].param.options, tests[i].param.correct)
			if !tests[i].expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if tests[i].expectError && err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

type multipleParam struct {
	options []string
	correct []string
}

func TestNewQuestion(t *testing.T) {
	tests := []struct {
		name        string
		param       questionParam
		expectError bool
	}{
		{
			name: "text is empty",
			param: questionParam{
				text:    "",
				details: newValidDetails(),
				score:   1,
			},
			expectError: true,
		},
		{
			name: "details is nil",
			param: questionParam{
				text:    "text",
				details: nil,
				score:   1,
			},
			expectError: true,
		},
		{
			name: "score equal 0",
			param: questionParam{
				text:    "text",
				details: newValidDetails(),
				score:   0,
			},
			expectError: true,
		},
		{
			name: "score less 0",
			param: questionParam{
				text:    "text",
				details: newValidDetails(),
				score:   -10,
			},
			expectError: true,
		},
		{
			name: "details is invalid",
			param: questionParam{
				text:    "text",
				details: newInvalidDetails(),
				score:   1,
			},
			expectError: true,
		},
		{
			name: "question is valid",
			param: questionParam{
				text:    "text",
				details: newValidDetails(),
				score:   1,
			},
			expectError: false,
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			_, err := NewQuestion(
				tests[i].param.text,
				tests[i].param.details,
				tests[i].param.score,
			)

			if !tests[i].expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if tests[i].expectError && err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

type questionParam struct {
	text    string
	details QuestionDetails
	score   int
}

func newValidDetails() QuestionDetails {
	return &SingleChoiceQuestion{
		Options: []string{"opt1", "opt2"},
		Correct: "opt1",
	}
}

func newInvalidDetails() QuestionDetails {
	return &SingleChoiceQuestion{
		Options: nil,
		Correct: "opt1",
	}
}

func TestNewQuiz(t *testing.T) {
	tests := []struct {
		name        string
		param       quizParam
		expectError bool
	}{
		{
			name:        "ownerID is invalid",
			param:       newInvalidQuizParamOwnerID(t),
			expectError: true,
		},
		{
			name:        "subjectID is invalid",
			param:       newInvalidQuizParamSubjectID(t),
			expectError: true,
		},
		{
			name:        "title is empty",
			param:       newInvalidQuizParamTitle(t),
			expectError: true,
		},
		{
			name:        "summary is empty",
			param:       newInvalidQuizParamSummary(t),
			expectError: true,
		},
		{
			name:        "questions are empty",
			param:       newInvalidQuizParamQuestions(t),
			expectError: true,
		},
		{
			name:        "maxAttempts equal 0",
			param:       newQuizParamMaxAttempts(t, 0),
			expectError: true,
		},
		{
			name:        "maxAttempts less 0",
			param:       newQuizParamMaxAttempts(t, -10),
			expectError: true,
		},
		{
			name:        "deadline is nil",
			param:       newQuizParamDeadline(t, nil),
			expectError: false,
		},
		{
			name:        "expired deadline",
			param:       newQuizParamDeadline(t, newInvalidDeadline()),
			expectError: true,
		},
		{
			name:        "quiz is valid",
			param:       newValidQuizParam(t),
			expectError: false,
		},
	}

	for i := range tests {
		t.Run(tests[i].name, func(t *testing.T) {
			_, err := NewQuiz(
				tests[i].param.ownerID,
				tests[i].param.subjectID,
				tests[i].param.title,
				tests[i].param.summary,
				tests[i].param.questions,
				tests[i].param.maxAttempts,
				tests[i].param.deadline,
			)

			if !tests[i].expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if tests[i].expectError && err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

type quizParam struct {
	title, summary     string
	ownerID, subjectID uuid.UUID
	questions          []Question
	maxAttempts        int
	deadline           *time.Time
}

func newValidQuizParam(t *testing.T) quizParam {
	return quizParam{
		title:       "title",
		summary:     "summary",
		ownerID:     uuid.New(),
		subjectID:   uuid.New(),
		questions:   newValidQuestions(t),
		maxAttempts: 5,
		deadline:    nil,
	}
}

func newValidQuestions(t *testing.T) []Question {
	questions := make([]Question, 5)
	for i := range questions {
		question, err := NewQuestion("text", newValidDetails(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		questions[i] = question
	}
	return questions
}

func newInvalidQuizParamOwnerID(t *testing.T) quizParam {
	param := newValidQuizParam(t)
	param.ownerID = uuid.Nil
	return param
}

func newInvalidQuizParamSubjectID(t *testing.T) quizParam {
	param := newValidQuizParam(t)
	param.subjectID = uuid.Nil
	return param
}

func newInvalidQuizParamTitle(t *testing.T) quizParam {
	param := newValidQuizParam(t)
	param.title = ""
	return param
}

func newInvalidQuizParamSummary(t *testing.T) quizParam {
	param := newValidQuizParam(t)
	param.summary = ""
	return param
}

func newInvalidQuizParamQuestions(t *testing.T) quizParam {
	param := newValidQuizParam(t)
	param.questions = nil

	return param
}

func newQuizParamMaxAttempts(t *testing.T, maxAttempts int) quizParam {
	param := newValidQuizParam(t)
	param.maxAttempts = maxAttempts

	return param
}

func newQuizParamDeadline(t *testing.T, deadline *time.Time) quizParam {
	param := newValidQuizParam(t)
	param.deadline = deadline

	return param
}

func newInvalidDeadline() *time.Time {
	deadline := time.Now().UTC().Add(-time.Hour)
	return &deadline
}

func TestCheckAttempt(t *testing.T) {
	var (
		qID1 = uuid.New()
		qID2 = uuid.New()
		qID3 = uuid.New()
	)
	quiz := Quiz{
		Questions: []Question{
			{
				ID: qID1,
				Details: &SingleChoiceQuestion{
					Options: []string{"opt1", "opt2"},
					Correct: "opt1",
				},
				Score: 1,
			},
			{
				ID: qID2,
				Details: &MultipleChoiceQuestion{
					Options: []string{"opt1", "opt2"},
					Correct: []string{"opt1", "opt2"},
				},
				Score: 2,
			},
			{
				ID: qID3,
				Details: &NumericQuestion{
					Correct: 5,
				},
				Score: 3,
			},
		},
	}

	attempt := Attempt{
		Answers: []Answer{
			{
				QuestionID: qID1,
				Answer:     "opt1",
			},
			{
				QuestionID: qID2,
				Answer:     []string{"opt1", "opt2"},
			},
			{
				QuestionID: qID3,
				Answer:     5,
			},
		},
	}

	err := quiz.CheckAttempt(&attempt)
	assert.NoError(t, err)
	assert.Equal(t, attempt.Score, 6)
}
