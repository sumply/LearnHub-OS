package get_user

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Postgres interface {
	Users(context.Context) ([]domain.User, error)
	StudentMap(context.Context) (map[uuid.UUID]domain.Student, error)
	TeacherMap(context.Context) (map[uuid.UUID]domain.Teacher, error)
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
	users, err := u.postgres.Users(ctx)
	if err != nil {
		return Output{}, err
	}

	teachers, err := u.postgres.TeacherMap(ctx)
	if err != nil {
		return Output{}, err
	}

	students, err := u.postgres.StudentMap(ctx)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		Users: make([]OutputUser, len(users)),
	}
	for i := range users {
		output.Users[i] = OutputUser{
			ID:        users[i].ID,
			FirstName: users[i].FirstName,
			LastName:  users[i].LastName,
			Role:      users[i].Role.String(),
		}
		switch users[i].Role {
		case domain.RoleStudent:
			group := students[users[i].ID].Group
			output.Users[i].Group = &group
		case domain.RoleTeacher:
			output.Users[i].Groups = teachers[users[i].ID].Groups
		}
	}

	return output, nil
}
