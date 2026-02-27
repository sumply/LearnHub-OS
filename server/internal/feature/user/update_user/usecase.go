package update_user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type Postgres interface {
	UpdateProfile(context.Context, uuid.UUID, string, string) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (uc *UseCase) UpdateUser(ctx context.Context, identity usecase.Identity, input Input) error {
	if err := uc.validateAccess(identity, input); err != nil {
		return err
	}

	if err := uc.postgres.UpdateProfile(ctx, input.ID, input.FirstName, input.LastName); err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) validateAccess(identity usecase.Identity, input Input) error {
	if identity.Role() != domain.RoleAdmin {
		if identity.ID() != input.ID {
			return usecase.NewAuthError(
				fmt.Errorf("user has not access to update profile with id %s", input.ID),
			)
		}
	}
	return nil
}
