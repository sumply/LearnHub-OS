package delete_subject

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type Subject interface {
	repository.Remover[domain.Subject]
}

type UseCase struct {
	subject Subject
}

func New(subject Subject) *UseCase {
	return &UseCase{
		subject: subject,
	}
}

func (u *UseCase) DeleteSubject(ctx context.Context, identity usecase.Identity, id uuid.UUID) error {
	if identity.Role() != domain.RoleAdmin {
		return usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	return u.subject.Remove(ctx, id)
}
