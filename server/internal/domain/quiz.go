package domain

import (
	"fmt"
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

func NewQuiz(title, summary string, ownerID, subjectID uuid.UUID, content []QuestionAggregate, maxAttempts int, deadline *time.Time) (Quiz, error) {
	if len(content) == 0 {
		return Quiz{}, fmt.Errorf("content is empty: %w", ErrInvalid)
	}

	totalCount := 0
	for i := range content {
		totalCount += content[i].Score
	}

	if deadline != nil && deadline.Before(time.Now().UTC()) {
		return Quiz{}, fmt.Errorf("deadline is before now: %w", ErrInvalid)
	}

	if maxAttempts <= 0 {
		return Quiz{}, fmt.Errorf("max_attempts less 0: %w", ErrInvalid)
	}

	return Quiz{
		ID:          uuid.New(),
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

type QuestionAggregate struct {
	ID      uuid.UUID
	QuizID  uuid.UUID
	Type    QuestionType
	Text    string
	Details QuestionDetails
	Score   int
}

type QuestionDetails interface {
	CheckAnswer(any) (bool, error)
}

func (q *QuestionAggregate) Validate() error {
	switch q.Type {
	case TypeSingleChoice:
		return q.validateSingleChoiceType()
	case TypeMultipleChoice:
		return q.validateMultipleType()
	case TypeNumeric:
		return q.validateNumericType()
	default:
		return fmt.Errorf("question has an incorrect type (%v): %w", q.Type, ErrInvalid)
	}
}

func (q *QuestionAggregate) validateSingleChoiceType() error {
	question, ok := q.Details.(*SingleChoiceQuestion)
	if !ok {
		return fmt.Errorf("question has an incorrect structure (%T): %w", q.Details, ErrInvalid)
	}
	err := question.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (q *QuestionAggregate) validateNumericType() error {
	question, ok := q.Details.(*NumericQuestion)
	if !ok {
		return fmt.Errorf("question has an incorrect structure: %w", ErrInvalid)
	}
	err := question.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (q *QuestionAggregate) validateMultipleType() error {
	question, ok := q.Details.(*MultipleChoiceQuestion)
	if !ok {
		return fmt.Errorf("question has an incorrect structure: %w", ErrInvalid)
	}
	err := question.Validate()
	if err != nil {
		return err
	}
	return nil
}

type SingleChoiceQuestion struct {
	Options []string `json:"options"`
	Correct int      `json:"correct"`
}

func (s *SingleChoiceQuestion) CheckAnswer(answer any) (bool, error) {
	_answer, ok := answer.(int)
	if !ok {
		return false, fmt.Errorf("answer is incorrect type: %w", ErrValidate)
	}
	if s.Correct != _answer {
		return false, nil
	}
	return true, nil
}

func (s *SingleChoiceQuestion) Validate() error {
	if len(s.Options) == 0 {
		return fmt.Errorf("options is empty: %w", ErrInvalid)
	}

	if s.Correct < 0 || len(s.Options) <= s.Correct {
		return fmt.Errorf("correct index is out of range: %w", ErrInvalid)
	}

	return nil
}

type MultipleChoiceQuestion struct {
	Options []string `json:"options"`
	Correct []int    `json:"correct"`
}

func (m *MultipleChoiceQuestion) CheckAnswer(answer any) (bool, error) {
	_answer, ok := answer.([]int)
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

func (m *MultipleChoiceQuestion) Validate() error {
	if len(m.Options) == 0 {
		return fmt.Errorf("options is empty: %w", ErrInvalid)
	}

	if len(m.Correct) == 0 {
		return fmt.Errorf("correct is empty: %w", ErrInvalid)
	}

	if !m.correctContainsUnique() {
		return fmt.Errorf("correct contains double: %w", ErrInvalid)
	}

	err := m.validateCorrect()
	if err != nil {
		return err
	}

	return nil
}

func (m *MultipleChoiceQuestion) correctContainsUnique() bool {
	set := make(map[int]struct{})
	for i := range m.Correct {
		if _, ok := set[m.Correct[i]]; ok {
			return false
		}
		set[m.Correct[i]] = struct{}{}
	}
	return true
}

func (m *MultipleChoiceQuestion) validateCorrect() error {
	for correct := range m.Correct {
		if correct < 0 || len(m.Options) <= correct {
			return fmt.Errorf("correct index is out of range: %w", ErrInvalid)
		}
	}
	return nil
}

type NumericQuestion struct {
	Correct float64 `json:"correct"`
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

func (n *NumericQuestion) Validate() error {
	return nil
}

type QuestionType string

const (
	TypeSingleChoice   QuestionType = "single"
	TypeMultipleChoice QuestionType = "multiple"
	TypeNumeric        QuestionType = "numeric"
)
