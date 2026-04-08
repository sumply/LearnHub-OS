package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Attempt struct {
	*Postgres
}

func NewAttempt(p *Postgres) *Attempt {
	return &Attempt{
		Postgres: p,
	}
}

type answerVisitor struct {
	ctx context.Context
	q   *sqlc.Queries
}

func (a *answerVisitor) VisitSingleAnswer(answer *domain.SingleAnswer) error {
	if err := a.saveQuizAnswer(answer); err != nil {
		return err
	}
	return a.q.InsertQuizAnswerSingle(a.ctx, sqlc.InsertQuizAnswerSingleParams{
		AnswerID:       answer.ID(),
		SelectedAnswer: answer.Selected,
	})
}

func (a *answerVisitor) VisitMultipleAnswer(answer *domain.MultipleAnswer) error {
	if err := a.saveQuizAnswer(answer); err != nil {
		return err
	}
	return a.q.InsertQuizAnswerMultiple(a.ctx, sqlc.InsertQuizAnswerMultipleParams{
		AnswerID:       answer.ID(),
		SelectedAnswer: answer.Selected,
	})
}

func (a *answerVisitor) VisitNumericAnswer(answer *domain.NumericAnswer) error {
	if err := a.saveQuizAnswer(answer); err != nil {
		return err
	}
	return a.q.InsertQuizAnswerNumeric(a.ctx, sqlc.InsertQuizAnswerNumericParams{
		AnswerID:       answer.ID(),
		SelectedAnswer: float64(answer.Selected),
	})
}

func (a *answerVisitor) saveQuizAnswer(answer domain.IAnswer) error {
	return a.q.InsertQuizAnswer(a.ctx, sqlc.InsertQuizAnswerParams{
		ID:         answer.ID(),
		AttemptID:  answer.AttemptID(),
		QuestionID: answer.QuestionID(),
	})
}

func (a *Attempt) Save(ctx context.Context, attempt *domain.Attempt) error {
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

		saver := &answerVisitor{
			ctx: ctx,
			q:   q,
		}

		for _, answer := range attempt.SelectedAnswers {
			if err := answer.Accept(saver); err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *Attempt) Get(ctx context.Context, id uuid.UUID) (*domain.Attempt, error) {
	row, err := a.sqlc.GetDomainAttempt(ctx, id)
	if err != nil {
		return nil, err
	}

	var answerRows []jsonb.AnswerRow
	if err := json.Unmarshal(row.Answers, &answerRows); err != nil {
		return nil, err
	}

	return &domain.Attempt{
		ID:         row.ID,
		QuizID:     row.QuizID,
		UserID:     row.UserID,
		TotalScore: row.Score,
		StartedAt:  row.StartedAt,
		EndedAt: func() *time.Time {
			if row.EndedAt.Valid {
				return &row.EndedAt.Time
			}
			return nil
		}(),
		SelectedAnswers: func() []domain.IAnswer {
			answers := make([]domain.IAnswer, 0, len(answerRows))
			for _, row := range answerRows {
				answers = append(answers, row.Domain)
			}
			return answers
		}(),
	}, nil
}

func (a *Attempt) GetAll(ctx context.Context) ([]*domain.Attempt, error) {
	rows, err := a.sqlc.GetDomainAttempts(ctx)
	if err != nil {
		return nil, err
	}

	attempts := make([]*domain.Attempt, 0, len(rows))
	for _, row := range rows {
		attempts = append(attempts, &domain.Attempt{
			ID:         row.ID,
			QuizID:     row.QuizID,
			UserID:     row.UserID,
			TotalScore: row.Score,
			StartedAt:  row.StartedAt,
			EndedAt: func() *time.Time {
				if row.EndedAt.Valid {
					return &row.EndedAt.Time
				}
				return nil
			}(),
		})
	}

	return attempts, nil
}

type answerFinisher struct {
	ctx context.Context
	q   *sqlc.Queries
}

func (f *answerFinisher) VisitSingleAnswer(a *domain.SingleAnswer) error {
	if err := f.finishAnswer(a); err != nil {
		return err
	}
	return f.q.FinishSingleAnswer(f.ctx, sqlc.FinishSingleAnswerParams{
		SelectedAnswer: a.Selected,
		AnswerID:       a.ID(),
	})
}

func (f *answerFinisher) VisitMultipleAnswer(a *domain.MultipleAnswer) error {
	if err := f.finishAnswer(a); err != nil {
		return err
	}
	return f.q.FinishMultipleAnswer(f.ctx, sqlc.FinishMultipleAnswerParams{
		SelectedAnswer: a.Selected,
		AnswerID:       a.ID(),
	})
}

func (f *answerFinisher) VisitNumericAnswer(a *domain.NumericAnswer) error {
	if err := f.finishAnswer(a); err != nil {
		return err
	}
	return f.q.FinishNumericAnswer(f.ctx, sqlc.FinishNumericAnswerParams{
		SelectedAnswer: float64(a.Selected),
		AnswerID:       a.ID(),
	})
}

func (f *answerFinisher) finishAnswer(a domain.IAnswer) error {
	return f.q.FinishAnswer(f.ctx, sqlc.FinishAnswerParams{
		Score:     a.Score(),
		IsCorrect: a.IsCorrect(),
		ID:        a.ID(),
	})
}

func (a *Attempt) Update(ctx context.Context, attempt *domain.Attempt) error {
	return a.withQueries(ctx, func(q *sqlc.Queries) error {
		err := a.sqlc.FinishAttempt(ctx, sqlc.FinishAttemptParams{
			Score: attempt.TotalScore,
			EndedAt: func() sql.NullTime {
				if attempt.EndedAt == nil {
					return sql.NullTime{}
				}
				return sql.NullTime{
					Time:  *attempt.EndedAt,
					Valid: true,
				}
			}(),
			ID: attempt.ID,
		})
		if err != nil {
			return err
		}

		visitor := &answerFinisher{
			ctx: ctx,
			q:   q,
		}

		for _, answer := range attempt.SelectedAnswers {
			if err := answer.Accept(visitor); err != nil {
				return err
			}
		}

		return nil
	})
}
