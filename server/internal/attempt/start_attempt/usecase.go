package start_attempt

import (
	"context"
	"server/internal/domain"
	"server/internal/query"

	"github.com/google/uuid"
)

type Postgres interface {
	DomainQuiz(context.Context, uuid.UUID) (domain.Quiz, error)
	QuizWithoutAnswers(context.Context, uuid.UUID) (query.Quiz, error)
}

type Redis interface {
	SaveAttempt(context.Context, *domain.Attempt) error
}

type UseCase struct {
	postgres Postgres
	redis    Redis
}

func New(postgres Postgres, redis Redis) *UseCase {
	return &UseCase{
		postgres: postgres,
		redis:    redis,
	}
}

func (u *UseCase) StartAttempt(ctx context.Context, input *Input) (Output, error) {
	quiz, err := u.postgres.DomainQuiz(ctx, input.QuizID)
	if err != nil {
		return Output{}, err
	}

	attempt, err := domain.NewAttempt(&quiz, input.UserID)
	if err != nil {
		return Output{}, err
	}

	err = u.redis.SaveAttempt(ctx, &attempt)
	if err != nil {
		return Output{}, err
	}

	outputQuiz, err := u.postgres.QuizWithoutAnswers(ctx, input.QuizID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Attempt: OutputAttempt{
			ID:        attempt.ID,
			StartedAt: attempt.StartedAt,
		},
		Quiz: outputQuiz,
	}, nil
}
