package get_user

import (
	"context"
	"server/internal/domain"
)

type Postgres interface {
	DetailedUsers(context.Context) ([]domain.DetailedUser, error)
}

type UseCase struct {
	postgres Postgres
}

func New(postgres Postgres) *UseCase {
	return &UseCase{
		postgres: postgres,
	}
}

func (u *UseCase) GetUser(ctx context.Context) (Output, error) {
	users, err := u.postgres.DetailedUsers(ctx)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Users: make([]OutputUser, len(users)),
	}
	for i := range users {
		switch u := users[i].(type) {
		case *domain.User:
			output.Users[i] = OutputUser{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Role:      u.Role.String(),
			}
		case *domain.Student:
			output.Users[i] = OutputUser{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Role:      u.Role.String(),
				Group:     &u.Group,
			}
		case *domain.Teacher:
			output.Users[i] = OutputUser{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Role:      u.Role.String(),
				Groups:    u.Groups,
				Subjects:  u.Subjects,
			}
		}
	}

	return output, nil
}
