package get_by_id

import (
	"context"
	"server/internal/query"

	"github.com/google/uuid"
)

type Postgres interface {
	Quiz(context.Context, uuid.UUID) (query.Quiz, error)
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
	quiz, err := u.postgres.Quiz(ctx, id)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Quiz: quiz,
	}
	return output, nil
}
