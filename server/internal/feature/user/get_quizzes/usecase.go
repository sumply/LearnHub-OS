package get_quizzes

import (
	"context"
	"server/internal/dto"

	"github.com/google/uuid"
)

type Repository interface {
	FindQuizLastAttempt(ctx context.Context, userID uuid.UUID) ([]dto.QuizLastAttempt, error)
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) GetQuizzes(ctx context.Context, userID uuid.UUID) (Output, error) {
	quizzes, err := u.repository.FindQuizLastAttempt(ctx, userID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Quizzes: quizzes,
	}, nil
}
