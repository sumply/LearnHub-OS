package get_users

import (
	"context"
	"server/internal/dto"

	"github.com/google/uuid"
)

type Repository interface {
	FindUserLastAttempt(context.Context, uuid.UUID) ([]dto.UserLastAttempt, error)
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) GetUsers(ctx context.Context, quizID uuid.UUID) (Output, error) {
	users, err := u.repository.FindUserLastAttempt(ctx, quizID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Users: users,
	}, nil
}
