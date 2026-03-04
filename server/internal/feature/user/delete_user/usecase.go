package delete_user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type User interface {
	repository.Remover[domain.User]
}

type UseCase struct {
	user User
}

func New(user User) *UseCase {
	return &UseCase{
		user: user,
	}
}

func (u *UseCase) DeleteUser(ctx context.Context, identity usecase.Identity, id uuid.UUID) error {
	if identity.Role() != domain.RoleAdmin {
		return usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	return u.user.Remove(ctx, id)
}
