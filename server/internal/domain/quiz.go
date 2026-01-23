package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID         uuid.UUID           `db:"id"`
	Title      string              `db:"title"`
	Summary    string              `db:"summary"`
	OwnerID    uuid.UUID           `db:"owner_id"`
	SubjectID  uuid.UUID           `db:"subject_id"`
	Content    []QuestionAggregate `db:"content"`
	TotalScore int                 `db:"total_score"`
	CreatedAt  time.Time           `db:"created_at"`
}

type QuestionAggregate struct {
	ID      uuid.UUID
	Type    QuestionType
	Text    string
	Payload any
	Score   int
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
		return fmt.Errorf("question has an incorrect type: %w", ErrInvalid)
	}
}

func (q *QuestionAggregate) validateSingleChoiceType() error {
	question, ok := q.Payload.(*SingleChoiceQuestion)
	if !ok {
		return fmt.Errorf("question has an incorrect structure: %w", ErrInvalid)
	}
	err := question.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (q *QuestionAggregate) validateNumericType() error {
	question, ok := q.Payload.(*NumericQuestion)
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
	question, ok := q.Payload.(*MultipleChoiceQuestion)
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
	Options []string
	Correct int
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
	Options []string
	Correct []int
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
	Tolerance float64
}

func (n *NumericQuestion) Validate() error {
	if n.Tolerance < 0 {
		return fmt.Errorf("tolerance less 0: %w", ErrInvalid)
	}

	return nil
}

type QuestionType string

const (
	TypeSingleChoice   QuestionType = "single"
	TypeMultipleChoice QuestionType = "multiple"
	TypeNumeric        QuestionType = "numeric"
)
