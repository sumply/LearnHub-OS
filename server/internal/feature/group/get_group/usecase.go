package get_group

import (
	"context"
	"server/internal/dto"
)

type QueryRepository interface {
	ListGroup(context.Context) ([]dto.Group, error)
}

type UseCase struct {
	qRepository QueryRepository
}

func New(qRepository QueryRepository) *UseCase {
	return &UseCase{
		qRepository: qRepository,
	}
}

func (u *UseCase) GetGroup(ctx context.Context) (Output, error) {
	groups, err := u.qRepository.ListGroup(ctx)
	if err != nil {
		return Output{}, err
	}

	return Output{
		Groups: groups,
	}, nil
}
