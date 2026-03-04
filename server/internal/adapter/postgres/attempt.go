package postgres

import (
	"context"
	"encoding/json"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Attempt struct {
	*Postgres
}

func NewAttemt(p *Postgres) *Attempt {
	return &Attempt{
		Postgres: p,
	}
}

func (a *Attempt) Save(ctx context.Context, attempt domain.Attempt) error {
	return a.withQueries(ctx, func(q *sqlc.Queries) error {
		err := a.sqlc.InsertQuizAttempt(ctx, sqlc.InsertQuizAttemptParams{
			ID:        attempt.ID,
			QuizID:    attempt.QuizID,
			UserID:    attempt.UserID,
			StartedAt: attempt.StartedAt,
		})
		if err != nil {
			return err
		}

		for _, answer := range attempt.Answers {
			details, err := json.Marshal(jsonb.AnswerDetails{Answer: answer.Answer})
			if err != nil {
				return err
			}

			err = a.sqlc.InsertQuizAnswer(ctx, sqlc.InsertQuizAnswerParams{
				ID:         answer.ID,
				AttemptID:  attempt.ID,
				QuestionID: answer.QuestionID,
				Details:    details,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *Attempt) Get(ctx context.Context, id uuid.UUID) (domain.Attempt, error) {
	return domain.Attempt{}, nil
}
