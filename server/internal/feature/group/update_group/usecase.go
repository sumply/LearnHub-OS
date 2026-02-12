package update_group

import (
	"context"

	"github.com/google/uuid"
)

type Postgres interface {
	UpdateGroup(context.Context, uuid.UUID, string) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) UpdateGroup(ctx context.Context, id uuid.UUID, input *Input) error {
	return u.postgres.UpdateGroup(ctx, id, input.Name)
}
