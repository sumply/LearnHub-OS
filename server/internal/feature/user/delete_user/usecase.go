package delete_user

import (
	"context"

	"github.com/google/uuid"
)

type Postgres interface {
	DeleteUser(context.Context, uuid.UUID) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.postgres.DeleteUser(ctx, id)
}
