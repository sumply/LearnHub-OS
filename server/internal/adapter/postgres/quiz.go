package postgres

import (
	"context"
	"database/sql"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Quiz struct {
	*Postgres
}

func NewQuiz(p *Postgres) *Quiz {
	return &Quiz{
		Postgres: p,
	}
}

func (q *Quiz) Save(ctx context.Context, quiz domain.Quiz) error {
	var deadline sql.NullTime
	if quiz.Deadline != nil {
		deadline = sql.NullTime{
			Time:  *quiz.Deadline,
			Valid: true,
		}
	}

	return q.withQueries(ctx, func(q *sqlc.Queries) error {
		err := q.InsertQuizInfo(ctx, sqlc.InsertQuizInfoParams{
			QuizID:      quiz.ID,
			Title:       quiz.Title,
			Summary:     quiz.Summary,
			SubjectID:   quiz.SubjectID,
			OwnerID:     quiz.OwnerID,
			MaxAttempts: int32(quiz.MaxAttempts),
			TotalScore:  quiz.TotalScore,
			Deadline:    deadline,
			CreatedAt:   quiz.CreatedAt,
		})
		if err != nil {
			return err
		}

		for _, id := range quiz.GroupIDs {
			err = q.InsertQuizAssigment(ctx, sqlc.InsertQuizAssigmentParams{
				QuizID:  quiz.ID,
				GroupID: id,
			})
			if err != nil {
				return err
			}
		}

		for _, question := range quiz.Questions {
			if err := question.Save(ctx, q); err != nil {
				return err
			}
		}
		return nil
	})
}

func (q *Quiz) Get(ctx context.Context, id uuid.UUID) (domain.Quiz, error) {

	return domain.Quiz{}, nil
}

func (q *Quiz) Remove(ctx context.Context, id uuid.UUID) error {
	return q.sqlc.DeleteQuizInfo(ctx, id)
}
