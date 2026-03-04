package get_by_id

import (
	"context"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type User interface {
	repository.Geter[domain.User]
}

type UseCase struct {
	User User
}

func New(user User) *UseCase {
	return &UseCase{
		User: user,
	}
}

func (u *UseCase) GetByID(ctx context.Context, id uuid.UUID) (Output, error) {
	user, err := u.User.Get(ctx, id)
	if err != nil {
		return Output{}, err
	}
	return Output{
		User: dto.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
		},
	}, nil
}
