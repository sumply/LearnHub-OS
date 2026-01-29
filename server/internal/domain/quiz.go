package domain

import (
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID          uuid.UUID
	Title       string
	Summary     string
	OwnerID     uuid.UUID
	SubjectID   uuid.UUID
	Content     []QuestionAggregate
	Deadline    *time.Time
	MaxAttempts int
	TotalScore  int
	CreatedAt   time.Time
}

func NewQuiz(id uuid.UUID, ownerID, subjectID uuid.UUID, title, summary string, content []QuestionAggregate, maxAttempts int, deadline *time.Time) (Quiz, error) {
	if id == uuid.Nil {
		return Quiz{}, fmt.Errorf("id is empty: %w", ErrValidate)
	}

	if ownerID == uuid.Nil {
		return Quiz{}, fmt.Errorf("ownerID is empty: %w", ErrValidate)
	}

	if subjectID == uuid.Nil {
		return Quiz{}, fmt.Errorf("ownerID is empty: %w", ErrValidate)
	}

	if title == "" {
		return Quiz{}, fmt.Errorf("title is empty: %w", ErrValidate)
	}

	if summary == "" {
		return Quiz{}, fmt.Errorf("summary is empty: %w", ErrValidate)
	}

	if len(content) == 0 {
		return Quiz{}, fmt.Errorf("content is empty: %w", ErrValidate)
	}

	totalCount := 0
	for i := range content {
		totalCount += content[i].Score
	}

	if maxAttempts <= 0 {
		return Quiz{}, fmt.Errorf("max_attempts less 0: %w", ErrValidate)
	}

	if deadline != nil && deadline.Before(time.Now().UTC()) {
		return Quiz{}, fmt.Errorf("deadline is before now: %w", ErrValidate)
	}

	return Quiz{
		ID:          id,
		Title:       title,
		Summary:     summary,
		OwnerID:     ownerID,
		SubjectID:   subjectID,
		Content:     content,
		Deadline:    deadline,
		MaxAttempts: maxAttempts,
		TotalScore:  totalCount,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func RestoreQuiz(
	id uuid.UUID,
	title, summary string,
	ownerID, subjectID uuid.UUID,
	content []QuestionAggregate,
	deadline *time.Time,
	maxAttempts, totalScore int,
	createdAt time.Time,
) (Quiz, error) {
	if id == uuid.Nil {
		return Quiz{}, fmt.Errorf("id is empty: %w", ErrInvalid)
	}
	if len(content) == 0 {
		return Quiz{}, fmt.Errorf("content is empty: %w", ErrInvalid)
	}
	if maxAttempts <= 0 {
		return Quiz{}, fmt.Errorf("max attempts less 0: %w", ErrInvalid)
	}
	if totalScore <= 0 {
		return Quiz{}, fmt.Errorf("total score less 0: %w", ErrInvalid)
	}
	return Quiz{
		ID:          id,
		Title:       title,
		Summary:     summary,
		OwnerID:     ownerID,
		SubjectID:   subjectID,
		Content:     content,
		Deadline:    deadline,
		MaxAttempts: maxAttempts,
		TotalScore:  totalScore,
		CreatedAt:   createdAt,
	}, nil
}

type QuestionAggregate struct {
	ID      uuid.UUID
	QuizID  uuid.UUID
	Type    QuestionType
	Text    string
	Details QuestionDetails
	Score   int
}

func newQuestionAggregate(id, quizID uuid.UUID, text string, details QuestionDetails, score int, errorType error) (QuestionAggregate, error) {
	if quizID == uuid.Nil {
		return QuestionAggregate{}, fmt.Errorf("quiz id is empty: %w", errorType)
	}

	if text == "" {
		return QuestionAggregate{}, fmt.Errorf("text is empty: %w", errorType)
	}

	if score <= 0 {
		return QuestionAggregate{}, fmt.Errorf("score less 0: %w", errorType)
	}

	if details == nil {
		return QuestionAggregate{}, fmt.Errorf("details is nil: %w", errorType)
	}

	var typ QuestionType
	switch details.(type) {
	case *SingleChoiceQuestion:
		typ = TypeSingleChoice
	case *MultipleChoiceQuestion:
		typ = TypeMultipleChoice
	case *NumericQuestion:
		typ = TypeNumeric
	default:
		return QuestionAggregate{}, fmt.Errorf("unknown question type: %w", errorType)
	}

	return QuestionAggregate{
		ID:      id,
		QuizID:  quizID,
		Type:    typ,
		Text:    text,
		Details: details,
		Score:   score,
	}, nil
}

func NewQuestionAggregate(quizID uuid.UUID, text string, details QuestionDetails, score int) (QuestionAggregate, error) {
	return newQuestionAggregate(uuid.New(), quizID, text, details, score, ErrValidate)
}

func RestoreQuestionAggregate(id, quizID uuid.UUID, text string, details QuestionDetails, score int, errorType error) (QuestionAggregate, error) {
	return newQuestionAggregate(id, quizID, text, details, score, ErrInvalid)
}

type QuestionDetails interface {
	CheckAnswer(answer any) (bool, error)
}

type SingleChoiceQuestion struct {
	Options []string `json:"options"`
	Correct string   `json:"correct"`
}

func newSingleChoiceQuestion(options []string, correct string, errorType error) (SingleChoiceQuestion, error) {
	if len(options) == 0 {
		return SingleChoiceQuestion{}, fmt.Errorf("options is nil: %w", errorType)
	}
	if !slices.Contains(options, correct) {
		return SingleChoiceQuestion{}, fmt.Errorf("options do not contain the correct: %w", errorType)
	}
	return SingleChoiceQuestion{
		Options: slices.Clone(options),
		Correct: correct,
	}, nil
}

func NewSingleChoiceQuestion(options []string, correct string) (SingleChoiceQuestion, error) {
	return newSingleChoiceQuestion(options, correct, ErrValidate)
}

func RestoreSingleChoiceQuestion(options []string, correct string) (SingleChoiceQuestion, error) {
	return newSingleChoiceQuestion(options, correct, ErrInvalid)
}

func (s *SingleChoiceQuestion) CheckAnswer(answer any) (bool, error) {
	_answer, ok := answer.(string)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type: %w", ErrValidate)
	}
	if s.Correct != _answer {
		return false, nil
	}
	return true, nil
}

type MultipleChoiceQuestion struct {
	Options []string `json:"options"`
	Correct []string `json:"correct"`
}

func newMultipleChoiceQuestion(options, correct []string, errorType error) (MultipleChoiceQuestion, error) {
	if len(options) == 0 {
		return MultipleChoiceQuestion{}, fmt.Errorf("options is empty: %w", errorType)
	}
	if len(correct) == 0 {
		return MultipleChoiceQuestion{}, fmt.Errorf("correct is empty: %w", errorType)
	}
	if len(correct) > len(options) {
		return MultipleChoiceQuestion{}, fmt.Errorf("more correct answers than options: %w", errorType)
	}
	for i := range correct {
		if !slices.Contains(options, correct[i]) {
			return MultipleChoiceQuestion{}, fmt.Errorf("options do not contain the correct: %w", errorType)
		}
	}
	return MultipleChoiceQuestion{
		Options: slices.Clone(options),
		Correct: slices.Clone(correct),
	}, nil
}

func NewMultipleChoiceQuestion(options, correct []string) (MultipleChoiceQuestion, error) {
	return newMultipleChoiceQuestion(options, correct, ErrValidate)
}

func RestoreMultipleChoiceQuestion(options, correct []string) (MultipleChoiceQuestion, error) {
	return newMultipleChoiceQuestion(options, correct, ErrInvalid)
}

func (m *MultipleChoiceQuestion) CheckAnswer(answer any) (bool, error) {
	_answer, ok := answer.([]string)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type: %w", ErrValidate)
	}
	for i := range m.Correct {
		if m.Correct[i] != _answer[i] {
			return false, nil
		}
	}
	return true, nil
}

type NumericQuestion struct {
	Correct float64 `json:"correct"`
}

func NewNumericQuestion(correct float64) NumericQuestion {
	return NumericQuestion{
		Correct: correct,
	}
}

func RestoreNumericQuestion(correct float64) NumericQuestion {
	return NumericQuestion{
		Correct: correct,
	}
}

func (n *NumericQuestion) CheckAnswer(answer any) (bool, error) {
	_answer, ok := answer.(float64)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type: %w", ErrValidate)
	}
	if n.Correct != _answer {
		return false, nil
	}
	return true, nil
}

type QuestionType string

const (
	TypeSingleChoice   QuestionType = "single"
	TypeMultipleChoice QuestionType = "multiple"
	TypeNumeric        QuestionType = "numeric"
)
