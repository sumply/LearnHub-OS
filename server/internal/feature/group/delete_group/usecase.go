package delete_group

import (
	"context"

	"github.com/google/uuid"
)

type Postgres interface {
	DeleteGroup(context.Context, uuid.UUID) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	return u.postgres.DeleteGroup(ctx, id)
}
