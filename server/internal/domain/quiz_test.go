package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewQuiz(t *testing.T) {
	tests := []struct {
		name        string
		param       newQuizParam
		expectError bool
	}{
		{
			name:        "invalid id",
			param:       newQuizParam{id: uuid.Nil},
			expectError: true,
		},
		{
			name: "invalid ownerID",
			param: newQuizParam{
				id:      uuid.New(),
				ownerID: uuid.Nil,
			},
			expectError: true,
		},
		{
			name: "invalid subjectID",
			param: newQuizParam{
				id:        uuid.New(),
				ownerID:   uuid.New(),
				subjectID: uuid.Nil,
			},
			expectError: true,
		},
		{
			name: "empty title",
			param: newQuizParam{
				id:        uuid.New(),
				ownerID:   uuid.New(),
				subjectID: uuid.New(),
				title:     "",
			},
			expectError: true,
		},
		{
			name: "empty summary",
			param: newQuizParam{
				id:        uuid.New(),
				ownerID:   uuid.New(),
				subjectID: uuid.New(),
				title:     "title",
				summary:   "",
			},
			expectError: true,
		},
		{
			name: "empty content",
			param: newQuizParam{
				id:        uuid.New(),
				ownerID:   uuid.New(),
				subjectID: uuid.New(),
				title:     "title",
				summary:   "summary",
				content:   nil,
			},
			expectError: true,
		},
		{
			name: "maxAttempts equal 0",
			param: newQuizParam{
				id:          uuid.New(),
				ownerID:     uuid.New(),
				subjectID:   uuid.New(),
				title:       "title",
				summary:     "summary",
				content:     []QuestionAggregate{QuestionAggregate{Score: 1}},
				maxAttempts: 0,
			},
			expectError: true,
		},
		{
			name: "maxAttempts less 0",
			param: newQuizParam{
				id:        uuid.New(),
				ownerID:   uuid.New(),
				subjectID: uuid.New(),
				title:     "title",
				summary:   "summary",
				content: []QuestionAggregate{
					{
						Score: 1,
					},
				},
				maxAttempts: -1,
			},
			expectError: true,
		},
		{
			name: "deadline is nil",
			param: newQuizParam{
				id:        uuid.New(),
				ownerID:   uuid.New(),
				subjectID: uuid.New(),
				title:     "title",
				summary:   "summary",
				content: []QuestionAggregate{
					{
						Score: 1,
					},
				},
				maxAttempts: 1,
				deadline:    nil,
			},
			expectError: false,
		},
		{
			name: "expired deadline",
			param: newQuizParam{
				id:        uuid.New(),
				ownerID:   uuid.New(),
				subjectID: uuid.New(),
				title:     "title",
				summary:   "summary",
				content: []QuestionAggregate{
					{
						Score: 1,
					},
				},
				maxAttempts: 1,
				deadline:    newExpiredDeadline(),
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewQuiz(
				test.param.id,
				test.param.ownerID,
				test.param.subjectID,
				test.param.title,
				test.param.summary,
				test.param.content,
				test.param.maxAttempts,
				test.param.deadline,
			)
			if !test.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if test.expectError && err == nil {
				t.Error("expected error but got nil")
			}
		})
	}
}

func newExpiredDeadline() *time.Time {
	deadline := time.Now().Add(-time.Hour)
	return &deadline
}

type newQuizParam struct {
	id          uuid.UUID
	title       string
	summary     string
	ownerID     uuid.UUID
	subjectID   uuid.UUID
	content     []QuestionAggregate
	maxAttempts int
	deadline    *time.Time
}

func TestNewQuestionAggregate(t *testing.T) {
	validQuizID := uuid.New()
	validDetails := &NumericQuestion{Correct: 42}

	tests := []struct {
		name        string
		quizID      uuid.UUID
		text        string
		details     QuestionDetails
		score       int
		expectError bool
	}{
		{
			name:        "valid question",
			quizID:      validQuizID,
			text:        "What is 2+2?",
			details:     &NumericQuestion{Correct: 4},
			score:       10,
			expectError: false,
		},
		{
			name:        "empty quiz id",
			quizID:      uuid.Nil,
			text:        "Text",
			details:     validDetails,
			score:       1,
			expectError: true,
		},
		{
			name:        "empty text",
			quizID:      validQuizID,
			text:        "",
			details:     validDetails,
			score:       1,
			expectError: true,
		},
		{
			name:        "score zero or less",
			quizID:      validQuizID,
			text:        "Text",
			details:     validDetails,
			score:       0,
			expectError: true,
		},
		{
			name:        "nil details",
			quizID:      validQuizID,
			text:        "Text",
			details:     nil,
			score:       1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewQuestionAggregate(tt.quizID, tt.text, tt.details, tt.score)
			if (err != nil) != tt.expectError {
				t.Errorf("NewQuestionAggregate() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// --- Тесты для SingleChoiceQuestion ---

func TestSingleChoiceQuestion_CheckAnswer(t *testing.T) {
	q, _ := NewSingleChoiceQuestion([]string{"A", "B", "C"}, "B")

	tests := []struct {
		name      string
		answer    any
		want      bool
		expectErr bool
	}{
		{"correct answer", "B", true, false},
		{"wrong answer", "A", false, false},
		{"wrong type", 123, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.CheckAnswer(tt.answer)
			if (err != nil) != tt.expectErr {
				t.Errorf("CheckAnswer() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAnswer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSingleChoiceQuestion_Validation(t *testing.T) {
	t.Run("correct not in options", func(t *testing.T) {
		_, err := NewSingleChoiceQuestion([]string{"A", "B"}, "C")
		if err == nil {
			t.Error("expected error because correct answer is not in options")
		}
	})
}

// --- Тесты для MultipleChoiceQuestion ---

func TestMultipleChoiceQuestion_CheckAnswer(t *testing.T) {
	q, _ := NewMultipleChoiceQuestion([]string{"A", "B", "C", "D"}, []string{"A", "C"})

	tests := []struct {
		name      string
		answer    any
		want      bool
		expectErr bool
	}{
		{"correct full answer", []string{"A", "C"}, true, false},
		{"partially correct", []string{"A", "B"}, false, false},
		{"wrong type", "A,C", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.CheckAnswer(tt.answer)
			if (err != nil) != tt.expectErr {
				t.Errorf("CheckAnswer() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAnswer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Тесты для NumericQuestion ---

func TestNumericQuestion_CheckAnswer(t *testing.T) {
	q := NewNumericQuestion(3.14)

	tests := []struct {
		name      string
		answer    any
		want      bool
		expectErr bool
	}{
		{"correct", 3.14, true, false},
		{"wrong", 3.15, false, false},
		{"wrong type", "3.14", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.CheckAnswer(tt.answer)
			if (err != nil) != tt.expectErr {
				t.Errorf("CheckAnswer() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAnswer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- Тесты для Restore функций ---

func TestRestoreQuiz(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	t.Run("successful restore", func(t *testing.T) {
		content := []QuestionAggregate{{Score: 10}}
		_, err := RestoreQuiz(id, "Title", "Sum", uuid.New(), uuid.New(), content, nil, 3, 10, now)
		if err != nil {
			t.Errorf("unexpected error on restore: %v", err)
		}
	})

	t.Run("invalid total score on restore", func(t *testing.T) {
		content := []QuestionAggregate{{Score: 10}}
		_, err := RestoreQuiz(id, "Title", "Sum", uuid.New(), uuid.New(), content, nil, 3, 0, now)
		if err == nil {
			t.Error("expected error due to totalScore <= 0")
		}
	})
}
