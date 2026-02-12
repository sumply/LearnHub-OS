package create_group

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Postgres interface {
	CreateGroup(context.Context, *domain.Group) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) CreateGroup(ctx context.Context, input *Input) (Output, error) {
	group := domain.Group{
		ID:   uuid.New(),
		Name: input.Name,
	}

	err := u.postgres.CreateGroup(ctx, &group)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: group.ID}, nil
}
