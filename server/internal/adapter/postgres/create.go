package postgres

import (
	"context"
	"encoding/json"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
)

func (p *Postgres) CreateGroup(ctx context.Context, group *domain.Group) error {
	err := p.sqlc.InsertSchoolGroup(ctx, sqlc.InsertSchoolGroupParams{
		ID:   group.ID,
		Name: group.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateSubject(ctx context.Context, subject *domain.Subject) error {
	err := p.sqlc.InsertSchoolSubject(ctx, sqlc.InsertSchoolSubjectParams{
		ID:   subject.ID,
		Name: subject.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateAttempt(ctx context.Context, attempt *domain.Attempt) error {
	err := p.sqlc.InsertQuizAttempt(ctx, sqlc.InsertQuizAttemptParams{
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

		err = p.sqlc.InsertQuizAnswer(ctx, sqlc.InsertQuizAnswerParams{
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
}
