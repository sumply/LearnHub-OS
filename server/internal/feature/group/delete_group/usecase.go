package delete_group

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type Group interface {
	repository.Remover[domain.Group]
}

type UseCase struct {
	group Group
}

func New(group Group) *UseCase {
	return &UseCase{
		group: group,
	}
}

func (u *UseCase) DeleteGroup(ctx context.Context, identity usecase.Identity, id uuid.UUID) error {
	if identity.Role() != domain.RoleAdmin {
		return usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	return u.group.Remove(ctx, id)
}
