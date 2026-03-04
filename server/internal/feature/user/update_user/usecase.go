package update_user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"
)

type User interface {
	repository.Geter[domain.User]
	repository.Updater[domain.User]
}

type UseCase struct {
	user User
}

func New(user User) *UseCase {
	return &UseCase{
		user: user,
	}
}

func (uc *UseCase) UpdateUser(ctx context.Context, identity usecase.Identity, input Input) error {
	if err := uc.validateAccess(identity, input); err != nil {
		return err
	}

	user, err := uc.user.Get(ctx, input.ID)
	if err != nil {
		return err
	}

	user.FirstName = input.FirstName
	user.LastName = input.LastName

	if err := uc.user.Update(ctx, user); err != nil {
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
