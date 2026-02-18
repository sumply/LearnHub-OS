package get_by_id

import (
	"context"
	"server/internal/dto"

	"github.com/google/uuid"
)

type Repository interface {
	User(context.Context, uuid.UUID) (dto.User, error)
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) GetByID(ctx context.Context, userID uuid.UUID) (Output, error) {
	user, err := u.repository.User(ctx, userID)
	if err != nil {
		return Output{}, err
	}
	return Output{
		User: user,
	}, nil
}
