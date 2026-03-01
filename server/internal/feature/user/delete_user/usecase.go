package delete_user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type DomainRepository interface {
	DeleteUser(context.Context, uuid.UUID) error
}

type UseCase struct {
	dRepository DomainRepository
}

func New(dRepository DomainRepository) *UseCase {
	return &UseCase{
		dRepository: dRepository,
	}
}

func (u *UseCase) DeleteUser(ctx context.Context, identity usecase.Identity, id uuid.UUID) error {
	if identity.Role() != domain.RoleAdmin {
		return usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	return u.dRepository.DeleteUser(ctx, id)
}
