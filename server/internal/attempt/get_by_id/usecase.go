package get_by_id

import (
	"context"
	"server/internal/dto"

	"github.com/google/uuid"
)

type Postgres interface {
	FinishedAttempt(context.Context, uuid.UUID) (dto.FinishedAttempt, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) GetByID(ctx context.Context, id uuid.UUID) (Output, error) {
	finishedAttempt, err := u.postgres.FinishedAttempt(ctx, id)
	if err != nil {
		return Output{}, err
	}

	return Output{
		FinishedAttempt: finishedAttempt,
	}, nil
}
