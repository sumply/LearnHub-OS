package start_attempt

import (
	"context"
	"server/internal/domain"
	"server/internal/dto"

	"github.com/google/uuid"
)

type Postgres interface {
	CreateAttempt(context.Context, *domain.Attempt) error
	DomainQuiz(context.Context, uuid.UUID) (domain.Quiz, error)
	Quiz(context.Context, uuid.UUID) (dto.Quiz, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
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

	outputQuiz, err := u.postgres.Quiz(ctx, input.QuizID)
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
