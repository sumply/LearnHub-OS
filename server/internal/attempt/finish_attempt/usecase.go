package finish_attempt

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Postgres interface {
	DomainQuestion(context.Context, uuid.UUID) (domain.Question, error)
	DomainAttempt(context.Context, uuid.UUID) (domain.Attempt, error)
	UpdateAttempt(context.Context, *domain.Attempt) error
}
type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) FinishAttempt(ctx context.Context, attemptID uuid.UUID, input *Input) (Output, error) {
	attempt, err := u.postgres.DomainAttempt(ctx, attemptID)
	if err != nil {
		return Output{}, err
	}

	for i, answer := range attempt.Answers {
		question, err := u.postgres.DomainQuestion(ctx, answer.QuestionID)
		if err != nil {
			return Output{}, err
		}
		ok, err := question.Details.CheckAnswer(answer.Answer)
		if err != nil {
			return Output{}, err
		}
		if ok {
			attempt.Score += answer.Score
		}
		answer.IsCorrect = ok
		attempt.Answers[i] = answer
	}

	err = u.postgres.UpdateAttempt(ctx, &attempt)
	if err != nil {
		return Output{}, err
	}

	return Output{}, nil
}
