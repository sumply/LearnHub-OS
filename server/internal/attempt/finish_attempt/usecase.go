package finish_attempt

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Postgres interface {
	UpdateAttempt(context.Context, *domain.Attempt) error
}

type Redis interface {
	Attempt(context.Context, uuid.UUID) (*domain.Attempt, error)
	Question(context.Context, uuid.UUID) (*domain.QuestionAggregate, error)
}

type UseCase struct {
	redis    Redis
	postgres Postgres
}

func New(postgres Postgres, redis Redis) *UseCase {
	return &UseCase{
		postgres: postgres,
		redis:    redis,
	}
}

func (u *UseCase) FinishAttempt(ctx context.Context, attemptID uuid.UUID, input *Input) (Output, error) {
	attempt, err := u.redis.Attempt(ctx, attemptID)
	if err != nil {
		return Output{}, err
	}

	for i, answer := range attempt.Answers {
		question, err := u.redis.Question(ctx, answer.QuestionID)
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

	err = u.postgres.UpdateAttempt(ctx, attempt)
	if err != nil {
		return Output{}, err
	}

	return Output{}, nil
}
