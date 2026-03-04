package login

import (
	"context"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/jwt"
	"server/internal/pkg/repository/filter"
)

type User interface {
	GetByCredential(context.Context, filter.Credential) (domain.User, error)
}

type UseCase struct {
	generator *jwt.Generator
	user      User
}

func New(generator *jwt.Generator, user User) *UseCase {
	return &UseCase{
		generator: generator,
		user:      user,
	}
}

func (u *UseCase) Login(ctx context.Context, input Input) (Output, error) {
	user, err := u.user.GetByCredential(ctx, filter.Credential{
		Email:   input.Email,
		PwdHash: "hash",
	})

	if err != nil {
		return Output{}, err
	}

	accessToken, err := u.generator.GenerateAccess(user.ID, user.Role)
	if err != nil {
		return Output{}, err
	}
	refreshToken, err := u.generator.GenerateRefresh(user.ID, user.Role)
	if err != nil {
		return Output{}, err
	}

	return Output{
		JWT: OutputJWT{
			Access:  accessToken,
			Refresh: refreshToken,
		},
		User: dto.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
		},
	}, nil
}
