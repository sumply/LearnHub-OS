package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"server/internal/adapter/postgres/jsonb"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/pkg/encoder"
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	*Postgres
}

type questionVisitor struct {
	ctx context.Context
	q   *sqlc.Queries
}

func (v *questionVisitor) VisitSingleQuestion(d *domain.SingleQuestion) error {
	err := v.q.InsertQuizQuestion(v.ctx, sqlc.InsertQuizQuestionParams{
		ID:     d.ID(),
		QuizID: d.QuizID,
		Title:  d.Title,
		Score:  d.Score(),
		Type:   sqlc.QuizQuestionTypeSingle,
	})
	if err != nil {
		return err
	}

	return v.q.InsertQuizQuestionSingle(v.ctx, sqlc.InsertQuizQuestionSingleParams{
		QuestionID: d.ID(),
		Correct:    d.Correct,
		Options:    d.Options,
	})
}

func (v *questionVisitor) VisitNumericQuestion(d *domain.NumericQuestion) error {
	err := v.q.InsertQuizQuestion(v.ctx, sqlc.InsertQuizQuestionParams{
		ID:     d.ID(),
		QuizID: d.QuizID,
		Title:  d.Title,
		Score:  d.Score(),
		Type:   sqlc.QuizQuestionTypeNumeric,
	})
	if err != nil {
		return err
	}

	return v.q.InsertQuizQuestionNumeric(v.ctx, sqlc.InsertQuizQuestionNumericParams{
		QuestionID: d.ID(),
		Correct:    float64(d.Correct),
	})
}

func (v *questionVisitor) VisitMultipleQuestion(d *domain.MultipleQuestion) error {
	err := v.q.InsertQuizQuestion(v.ctx, sqlc.InsertQuizQuestionParams{
		ID:     d.ID(),
		QuizID: d.QuizID,
		Title:  d.Title,
		Score:  d.Score(),
		Type:   sqlc.QuizQuestionTypeMultiple,
	})
	if err != nil {
		return err
	}

	return v.q.InsertQuizQuestionMultiple(v.ctx, sqlc.InsertQuizQuestionMultipleParams{
		QuestionID: d.ID(),
		Correct:    d.Correct,
		Options:    d.Options,
	})
}

func NewQuiz(p *Postgres) *Quiz {
	return &Quiz{
		Postgres: p,
	}
}

func (q *Quiz) Save(ctx context.Context, quiz *domain.Quiz) error {
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
			visitor := &questionVisitor{
				ctx: ctx,
				q:   q,
			}
			if err := question.Accept(visitor); err != nil {
				return err
			}
		}

		return nil
	})
}

func (q *Quiz) Get(ctx context.Context, id uuid.UUID) (*domain.Quiz, error) {
	row, err := q.sqlc.GetDomainQuiz(ctx, id)
	if err != nil {
		return nil, err
	}

	var questionRows []jsonb.QuestionRow
	if err := json.Unmarshal(row.Questions, &questionRows); err != nil {
		return nil, err
	}

	questions := make([]domain.IQuestion, 0, len(questionRows))
	for _, row := range questionRows {
		question, err := row.Unpack()
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	quiz := &domain.Quiz{
		ID:        row.QuizID,
		Title:     row.Title,
		Summary:   row.Summary,
		OwnerID:   row.OwnerID,
		SubjectID: row.SubjectID,
		GroupIDs:  row.GroupIds,
		Deadline: func() *time.Time {
			if row.Deadline.Valid {
				return &row.Deadline.Time
			}
			return nil
		}(),
		MaxAttempts: int(row.MaxAttempts),
		TotalScore:  row.TotalScore,
		CreatedAt:   row.CreatedAt,
		Questions:   questions,
	}

	encoder.JSON(os.Stdout, quiz)

	return quiz, nil
}

func (q *Quiz) GetAll(ctx context.Context) ([]*domain.Quiz, error) {
	rows, err := q.sqlc.GetDomainQuizzes(ctx)
	if err != nil {
		return nil, err
	}

	domains := make([]*domain.Quiz, 0, len(rows))
	for _, row := range rows {
		domains = append(domains, &domain.Quiz{
			ID:        row.QuizID,
			Title:     row.Title,
			Summary:   row.Summary,
			OwnerID:   row.OwnerID,
			SubjectID: row.SubjectID,
			GroupIDs:  row.GroupIds,
			Deadline: func() *time.Time {
				if row.Deadline.Valid {
					return &row.Deadline.Time
				}
				return nil
			}(),
			MaxAttempts: int(row.MaxAttempts),
			TotalScore:  row.TotalScore,
			CreatedAt:   row.CreatedAt,
		})
	}

	return domains, nil
}

func (q *Quiz) Remove(ctx context.Context, id uuid.UUID) error {
	return q.sqlc.DeleteQuizInfo(ctx, id)
}
