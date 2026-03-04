package update_group

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type Group interface {
	repository.Geter[domain.Group]
	repository.Updater[domain.Group]
}

type UseCase struct {
	group Group
}

func New(group Group) *UseCase {
	return &UseCase{
		group: group,
	}
}

func (u *UseCase) UpdateGroup(ctx context.Context, id uuid.UUID, input *Input) error {
	group, err := u.group.Get(ctx, id)
	if err != nil {
		return err
	}

	group.Name = input.Name

	if err := u.group.Update(ctx, group); err != nil {
		return err
	}

	return nil
}
