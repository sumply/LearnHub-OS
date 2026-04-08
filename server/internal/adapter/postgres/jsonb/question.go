package jsonb

import (
	"encoding/json"
	"fmt"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"

	"github.com/google/uuid"
)

type QuestionRow struct {
	domain domain.IQuestion
	err    error
}

func (q *QuestionRow) Unpack() (domain.IQuestion, error) {
	return q.domain, q.err
}

func (q *QuestionRow) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID              uuid.UUID             `json:"id"`
		Title           string                `json:"title"`
		Score           int                   `json:"score"`
		Type            sqlc.QuizQuestionType `json:"type"`
		SingleCorrect   string                `json:"single_correct"`
		SingleOptions   []string              `json:"single_options"`
		MultipleCorrect []string              `json:"multiple_correct"`
		MultipleOptions []string              `json:"multiple_options"`
		NumericCorrect  float32               `json:"numeric_correct"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	switch aux.Type {
	case sqlc.QuizQuestionTypeSingle:
		d, err := domain.RestoreSingleQuestion(
			aux.ID,
			aux.Title,
			aux.SingleCorrect,
			aux.SingleOptions,
			aux.Score,
		)
		q.domain = d
		q.err = err
	case sqlc.QuizQuestionTypeMultiple:
		d, err := domain.RestoreMultipleQuestion(
			aux.ID,
			aux.Title,
			aux.MultipleCorrect,
			aux.MultipleOptions,
			aux.Score,
		)
		q.domain = d
		q.err = err
	case sqlc.QuizQuestionTypeNumeric:
		d, err := domain.RestoreNumericQuestion(
			aux.ID,
			aux.Title,
			aux.NumericCorrect,
			aux.Score,
		)
		q.domain = d
		q.err = err
	default:
		return fmt.Errorf("invalid type")
	}
	return nil
}
