package create_user

import (
	"context"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Postgres interface {
	CreateUser(context.Context, *domain.User) error
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) CreateUser(ctx context.Context, input *Input) (Output, error) {
	role, err := domain.NewUserRole(input.Role)
	if err != nil {
		return Output{}, err
	}
	user := domain.User{
		ID:        uuid.New(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      role,
		Email:     input.Email,
		PwdHash:   "hash",
		CreatedAt: time.Now().UTC(),
	}

	if err := u.postgres.CreateUser(ctx, &user); err != nil {
		return Output{}, err
	}
	return Output{ID: user.ID}, nil
}
