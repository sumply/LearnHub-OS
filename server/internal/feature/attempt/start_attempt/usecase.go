package start_attempt

import (
	"context"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type Postgres interface {
	Quiz(context.Context, uuid.UUID) (dto.Quiz, error)
}

type Quiz interface {
	repository.Geter[domain.Quiz]
}

type Attempt interface {
	repository.Saver[domain.Attempt]
}

type UseCase struct {
	postgres Postgres
	quiz     Quiz
	attempt  Attempt
}

func New(postgres Postgres, quiz Quiz, attempt Attempt) *UseCase {
	return &UseCase{
		postgres: postgres,
		quiz:     quiz,
		attempt:  attempt,
	}
}

func (u *UseCase) StartAttempt(ctx context.Context, quizID uuid.UUID, input *Input) (Output, error) {
	quiz, err := u.quiz.Get(ctx, quizID)
	if err != nil {
		return Output{}, err
	}

	attempt, err := domain.NewAttempt(&quiz, input.UserID)
	if err != nil {
		return Output{}, err
	}

	err = u.attempt.Save(ctx, attempt)
	if err != nil {
		return Output{}, err
	}

	quizDTO, err := u.postgres.Quiz(ctx, quizID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Attempt: OutputAttempt{
			ID:        attempt.ID,
			StartedAt: attempt.StartedAt,
		},
		Quiz: quizDTO,
	}, nil
}
