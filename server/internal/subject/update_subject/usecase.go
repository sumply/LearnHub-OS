package update_subject

import (
	"context"

	"github.com/google/uuid"
)

type Postgres interface {
	UpdateSubject(context.Context, uuid.UUID, string) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) UpdateSubject(ctx context.Context, id uuid.UUID, input *Input) error {
	return u.postgres.UpdateSubject(ctx, id, input.Name)
}
