package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/adapter/postgres/jsonb"
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

	var questions jsonb.QuizQuestionAGGs
	if err := json.Unmarshal(row.QuizQuestions, &questions); err != nil {
		return domain.Quiz{}, err
	}

	quiz := domain.Quiz{
		ID:          row.QuizID,
		Title:       row.QuizTitle,
		Summary:     row.QuizSummary,
		OwnerID:     row.QuizOwnerID,
		SubjectID:   row.QuizSubjectID,
		Deadline:    deadline,
		Questions:   questions.ToDomainQuestions(),
		MaxAttempts: int(row.QuizMaxAttempts),
		TotalScore:  int(row.QuizTotalScore),
		CreatedAt:   row.QuizCreatedAt,
	}

	return quiz, nil
}

func (p *Postgres) DomainAttempt(ctx context.Context, id uuid.UUID) (domain.Attempt, error) {
	row, err := p.sqlc.GetDomainAttempt(ctx)
	if err != nil {
		return domain.Attempt{}, fmt.Errorf("getting domain attempt: %w", err)
	}

	var answers jsonb.QuizAnswerAGGs
	if err := json.Unmarshal(row.AttemptAnswers, &answers); err != nil {
		return domain.Attempt{}, err
	}

	var endedAt *time.Time
	if row.AttemptEndendAt.Valid {
		endedAt = &row.AttemptEndendAt.Time
	}

	return domain.Attempt{
		ID:        row.AttemptID,
		QuizID:    row.AttemptQuizID,
		UserID:    row.AttemptUserID,
		Answers:   answers.ToDomainAnswers(),
		Score:     int(row.AttemptScoreID),
		StartedAt: row.AttemptStartedAt,
		EndedAt:   endedAt,
	}, nil
}
