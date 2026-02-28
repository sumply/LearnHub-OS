package add_students

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type DomainRepository interface {
	AddStudentsToGroup(ctx context.Context, groupID uuid.UUID, studentIDs uuid.UUIDs) error
}

type UseCase struct {
	dRepository DomainRepository
}

func New(dRepository DomainRepository) *UseCase {
	return &UseCase{
		dRepository: dRepository,
	}
}

func (uc *UseCase) AddStudents(ctx context.Context, identity usecase.Identity, input Input) error {
	if identity.Role() != domain.RoleAdmin {
		return usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	return uc.dRepository.AddStudentsToGroup(ctx, input.GroupID, input.StudentIDs)
}
