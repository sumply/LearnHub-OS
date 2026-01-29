package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

func (p *Postgres) DomainQuiz(ctx context.Context, id uuid.UUID) (domain.Quiz, error) {
	row, err := p.sqlc.GetDomainQuiz(ctx, id)
	if err != nil {
		return domain.Quiz{}, fmt.Errorf("getting domain quiz: %w", err)
	}

	var deadline *time.Time
	if row.QuizDeadline.Valid {
		deadline = &row.QuizDeadline.Time
	}

	var content []domain.QuestionAggregate
	json.Unmarshal(row.QuizQuestions, &content)

	quiz := domain.Quiz{
		ID:          row.QuizID,
		Title:       row.QuizTitle,
		Summary:     row.QuizSummary,
		OwnerID:     row.QuizOwnerID,
		SubjectID:   row.QuizSubjectID,
		Deadline:    deadline,
		Content:     content,
		MaxAttempts: int(row.QuizMaxAttempts),
		TotalScore:  int(row.QuizTotalScore),
		CreatedAt:   row.QuizCreatedAt,
	}

	return quiz, nil
}

func (p *Postgres) DomainAttempt(ctx context.Context, id uuid.UUID) (domain.Attempt, error) {
	_, err := p.sqlc.GetDomainAttempt(ctx)
	if err != nil {
		return domain.Attempt{}, fmt.Errorf("getting domain attempt: %w", err)
	}

	return domain.Attempt{}, nil
}

func (p *Postgres) DomainQuestion(ctx context.Context, id uuid.UUID) (domain.QuestionAggregate, error) {
	return domain.QuestionAggregate{}, nil
}
