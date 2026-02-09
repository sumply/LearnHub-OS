package update_user

import (
	"context"

	"github.com/google/uuid"
)

type Postgres interface {
	UpdateProfile(context.Context, uuid.UUID, string, string) error
}

type Usecase struct {
	postgres Postgres
}

func New(postgres Postgres) *Usecase {
	return &Usecase{
		postgres: postgres,
	}
}

func (u *Usecase) UpdateUser(ctx context.Context, id uuid.UUID, input *Input) error {
	if err := u.postgres.UpdateProfile(ctx, id, input.FirstName, input.LastName); err != nil {
		return err
	}

	return nil
}
