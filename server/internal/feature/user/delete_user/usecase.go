package delete_user

import (
	"context"
	"server/internal/pkg/usecase"

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

func (u *UseCase) DeleteUser(ctx context.Context, identity usecase.Identity, id uuid.UUID) error {
	return u.postgres.DeleteUser(ctx, id)
}
