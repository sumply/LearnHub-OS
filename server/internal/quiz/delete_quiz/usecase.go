package delete_quiz

import (
	"context"

	"github.com/google/uuid"
)

type Postgres interface {
	DeleteQuiz(context.Context, uuid.UUID) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	err := u.postgres.DeleteQuiz(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
