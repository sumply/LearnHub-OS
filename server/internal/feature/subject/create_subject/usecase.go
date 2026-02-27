package create_subject

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/usecase"
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

func (u *UseCase) CreateSubject(ctx context.Context, identity usecase.Identity, input Input) (Output, error) {
	if identity.Role() != domain.RoleAdmin {
		return Output{}, usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	subject, err := domain.NewSubject(input.Name)
	if err != nil {
		return Output{}, err
	}

	if err := u.postgres.CreateSubject(ctx, &subject); err != nil {
		return Output{}, err
	}

	return Output{ID: subject.ID}, nil
}
