package get_students

import (
	"context"
	"server/internal/dto"

	"github.com/google/uuid"
)

type Repository interface {
	FindStudents(ctx context.Context, groupID uuid.UUID) ([]dto.User, error)
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) GetStudents(ctx context.Context, groupID uuid.UUID) (Output, error) {
	users, err := u.repository.FindStudents(ctx, groupID)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Users: users,
	}, nil
}
