package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
)

func (p *Postgres) UpdateAttempt(ctx context.Context, attempt *domain.Attempt) error {
	var endedAt sql.NullTime
	if attempt.EndedAt != nil {
		endedAt.Time = *attempt.EndedAt
		endedAt.Valid = true
	}
	err := p.sqlc.UpdateQuizAttempt(ctx, sqlc.UpdateQuizAttemptParams{
		Score:   attempt.Score,
		EndedAt: endedAt,
		ID:      attempt.ID,
	})
	if err != nil {
		return err
	}

	for _, answer := range attempt.Answers {
		details, err := json.Marshal(jsonb.AnswerDetails{Answer: answer.Answer})
		if err != nil {
			return err
		}
		err = p.sqlc.UpdateQuizAnswer(ctx, sqlc.UpdateQuizAnswerParams{
			Details:   details,
			Score:     answer.Score,
			IsCorrect: answer.IsCorrect,
			ID:        answer.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
