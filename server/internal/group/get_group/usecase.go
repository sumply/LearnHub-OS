package get_group

import (
	"context"
	"server/internal/domain"
)

type Postgres interface {
	Groups(context.Context) ([]domain.Group, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) GetGroup(ctx context.Context) (Output, error) {
	groups, err := u.postgres.Groups(ctx)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Groups: make([]OutputGroup, len(groups)),
	}
	for i := range groups {
		output.Groups[i] = OutputGroup{
			ID:   groups[i].ID,
			Name: groups[i].Name,
		}
	}

	return output, nil
}
