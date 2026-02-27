package delete_subject

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/usecase"

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

func (u *UseCase) DeleteSubject(ctx context.Context, identity usecase.Identity, id uuid.UUID) error {
	if identity.Role() != domain.RoleAdmin {
		return usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	return u.postgres.DeleteSubject(ctx, id)
}
