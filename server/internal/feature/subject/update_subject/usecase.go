package update_subject

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type Subject interface {
	repository.Geter[domain.Subject]
	repository.Updater[domain.Subject]
}

type UseCase struct {
	subject Subject
}

func New(subject Subject) *UseCase {
	return &UseCase{
		subject: subject,
	}
}

func (u *UseCase) UpdateSubject(ctx context.Context, id uuid.UUID, input *Input) error {
	subject, err := u.subject.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := u.subject.Update(ctx, subject); err != nil {
		return err
	}
	return nil
}
