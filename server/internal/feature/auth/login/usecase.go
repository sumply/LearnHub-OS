package login

import (
	"context"
	"server/internal/domain"
	"server/internal/repository/filter"
)

type Repository interface {
	UserByCredential(context.Context, filter.Credential) (domain.User, error)
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{
		repository: repository,
	}
}

func (u *UseCase) Login(ctx context.Context, input Input) (Output, error) {
	user, err := u.repository.UserByCredential(ctx, filter.Credential{
		Email:   input.Email,
		PwdHash: "hash",
	})

	if err != nil {
		return Output{}, err
	}

	return Output{
		ID:   user.ID,
		Role: user.Role,
	}, nil
}
