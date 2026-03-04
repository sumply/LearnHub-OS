package create_subject

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"
)

type Subject interface {
	repository.Saver[domain.Subject]
}

type UseCase struct {
	subject Subject
}

func New(subject Subject) *UseCase {
	return &UseCase{
		subject: subject,
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

	if err := u.subject.Save(ctx, subject); err != nil {
		return Output{}, err
	}

	return Output{ID: subject.ID}, nil
}
