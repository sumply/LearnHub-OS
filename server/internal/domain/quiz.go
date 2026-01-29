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
	Questions   []Question
	Deadline    *time.Time
	MaxAttempts int
	TotalScore  int
	CreatedAt   time.Time
}

func NewQuiz(
	ownerID, subjectID uuid.UUID,
	title, summary string,
	questions []Question,
	maxAttempts int,
	deadline *time.Time,
) (Quiz, error) {
	quizID := uuid.New()

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

	if len(questions) == 0 {
		return Quiz{}, fmt.Errorf("content is empty: %w", ErrValidate)
	}

	totalCount := 0
	for i := range questions {
		totalCount += questions[i].Score
		questions[i].QuizID = quizID
	}

	if maxAttempts <= 0 {
		return Quiz{}, fmt.Errorf("max_attempts less 0: %w", ErrValidate)
	}

	if deadline != nil && deadline.Before(time.Now().UTC()) {
		return Quiz{}, fmt.Errorf("deadline is before now: %w", ErrValidate)
	}

	return Quiz{
		ID:          quizID,
		Title:       title,
		Summary:     summary,
		OwnerID:     ownerID,
		SubjectID:   subjectID,
		Questions:   questions,
		Deadline:    deadline,
		MaxAttempts: maxAttempts,
		TotalScore:  totalCount,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

type Question struct {
	ID      uuid.UUID
	QuizID  uuid.UUID
	Text    string
	Details QuestionDetails
	Score   int
}

func NewQuestion(text string, details QuestionDetails, score int) (Question, error) {
	if text == "" {
		return Question{}, fmt.Errorf("text is empty: %w", ErrInvalid)
	}

	if score <= 0 {
		return Question{}, fmt.Errorf("score less 0: %w", ErrInvalid)
	}

	if details == nil {
		return Question{}, fmt.Errorf("details is nil: %w", ErrInvalid)
	}

	if err := details.Validate(); err != nil {
		return Question{}, err
	}

	return Question{
		ID:      uuid.New(),
		Text:    text,
		Details: details,
		Score:   score,
	}, nil
}

type QuestionDetails interface {
	Validate() error
	Variant() QuestionType
	CheckAnswer(answer any) (bool, error)
}

type SingleChoiceQuestion struct {
	Options []string
	Correct string
}

func NewSingleChoiceQuestion(options []string, correct string) (SingleChoiceQuestion, error) {
	single := SingleChoiceQuestion{
		Options: slices.Clone(options),
		Correct: correct,
	}
	if err := single.Validate(); err != nil {
		return SingleChoiceQuestion{}, err
	}
	return single, nil
}

func (s *SingleChoiceQuestion) Validate() error {
	if len(s.Options) == 0 {
		return fmt.Errorf("options is nil: %w", ErrInvalid)
	}
	if !slices.Contains(s.Options, s.Correct) {
		return fmt.Errorf("options do not contain the correct: %w", ErrInvalid)
	}
	return nil
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

func (s *SingleChoiceQuestion) Variant() QuestionType {
	return TypeSingleChoice
}

type MultipleChoiceQuestion struct {
	Options []string `json:"options"`
	Correct []string `json:"correct"`
}

func NewMultipleChoiceQuestion(options, correct []string) (MultipleChoiceQuestion, error) {
	multiple := MultipleChoiceQuestion{
		Options: slices.Clone(options),
		Correct: slices.Clone(correct),
	}
	if err := multiple.Validate(); err != nil {
		return MultipleChoiceQuestion{}, err
	}
	return multiple, nil
}

func (m *MultipleChoiceQuestion) Validate() error {
	if len(m.Options) == 0 {
		return fmt.Errorf("options is empty: %w", ErrInvalid)
	}
	if len(m.Correct) == 0 {
		return fmt.Errorf("correct is empty: %w", ErrInvalid)
	}
	if len(m.Correct) > len(m.Options) {
		return fmt.Errorf("more correct answers than options: %w", ErrInvalid)
	}
	for i := range m.Correct {
		if !slices.Contains(m.Options, m.Correct[i]) {
			return fmt.Errorf("options do not contain the correct: %w", ErrInvalid)
		}
	}
	return nil
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

func (m *MultipleChoiceQuestion) Variant() QuestionType {
	return TypeMultipleChoice
}

type NumericQuestion struct {
	Correct float64 `json:"correct"`
}

func NewNumericQuestion(correct float64) (NumericQuestion, error) {
	numeric := NumericQuestion{
		Correct: correct,
	}
	if err := numeric.Validate(); err != nil {
		return NumericQuestion{}, err
	}
	return numeric, nil
}

func (n *NumericQuestion) Validate() error {
	return nil
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

func (n *NumericQuestion) Variant() QuestionType {
	return TypeNumeric
}

type QuestionType string

const (
	TypeSingleChoice   QuestionType = "single"
	TypeMultipleChoice QuestionType = "multiple"
	TypeNumeric        QuestionType = "numeric"
)
