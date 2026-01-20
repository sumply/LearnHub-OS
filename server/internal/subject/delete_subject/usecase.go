package delete_subject

import (
	"context"

	"github.com/google/uuid"
)

type Postgres interface {
	DeleteSubject(context.Context, uuid.UUID) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	return u.DeleteSubject(ctx, id)
}
