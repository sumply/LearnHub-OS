package create_subject

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Postgres interface {
	CreateSubject(context.Context, *domain.Subject) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) CreateSubject(ctx context.Context, input *Input) (Output, error) {
	subject := domain.Subject{
		ID:   uuid.New(),
		Name: input.Name,
	}

	if err := u.postgres.CreateSubject(ctx, &subject); err != nil {
		return Output{}, err
	}

	return Output{ID: subject.ID}, nil
}
