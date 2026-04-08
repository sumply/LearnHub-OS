package jsonb

import (
	"encoding/json"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"

	"github.com/google/uuid"
)

type AnswerRow struct {
	AttemptID uuid.UUID
	Domain    domain.IAnswer
}

func (r *AnswerRow) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID                     uuid.UUID             `json:"id"`
		QuestionID             uuid.UUID             `json:"question_id"`
		Score                  int                   `json:"score"`
		IsCorrect              bool                  `json:"is_correct"`
		Type                   sqlc.QuizQuestionType `json:"type"`
		SingleSelectedAnswer   string                `json:"single_selected_answer"`
		MultipleSelectedAnswer []string              `json:"multiple_selected_answer"`
		NumericSelectedAnswer  float32               `json:"numeric_selected_answer"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	switch aux.Type {
	case sqlc.QuizQuestionTypeSingle:
		r.Domain = &domain.SingleAnswer{
			Answer: domain.RestoreAnswer(
				aux.ID,
				r.AttemptID,
				aux.QuestionID,
				aux.Score,
				aux.IsCorrect,
			),
			Selected: aux.SingleSelectedAnswer,
		}
	case sqlc.QuizQuestionTypeMultiple:
		r.Domain = &domain.MultipleAnswer{
			Answer: domain.RestoreAnswer(
				aux.ID,
				r.AttemptID,
				aux.QuestionID,
				aux.Score,
				aux.IsCorrect,
			),
			Selected: aux.MultipleSelectedAnswer,
		}
	case sqlc.QuizQuestionTypeNumeric:
		r.Domain = &domain.NumericAnswer{
			Answer: domain.RestoreAnswer(
				aux.ID,
				r.AttemptID,
				aux.QuestionID,
				aux.Score,
				aux.IsCorrect,
			),
			Selected: aux.NumericSelectedAnswer,
		}
	}

	return nil
}
