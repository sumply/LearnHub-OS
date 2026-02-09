package create_user

import (
	"context"
	"server/internal/domain"
	"server/internal/usecase"
)

type Postgres interface {
	CreateUser(context.Context, *domain.User) error
	CreateStudent(context.Context, *domain.Student) error
	CreateTeacher(context.Context, *domain.Teacher) error
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
	user, err := domain.NewUser(
		input.FirstName,
		input.LastName,
		input.Email,
		"hash",
		domain.UserRole(input.Role),
	)
	if err != nil {
		return Output{}, usecase.NewValidationError(err)
	}

	switch domain.UserRole(input.Role) {

	case domain.RoleAdmin:
		err := u.postgres.CreateUser(ctx, &user)
		if err != nil {
			return Output{}, err
		}

	case domain.RoleStudent:
		student, err := domain.NewStudent(&user, input.GroupID)
		if err != nil {
			return Output{}, err
		}
		err = u.postgres.CreateStudent(ctx, &student)
		if err != nil {
			return Output{}, err
		}

	case domain.RoleTeacher:
		teacher, err := domain.NewTeacher(&user, input.SubjectIDs, input.GroupIDs)
		if err != nil {
			return Output{}, err
		}
		err = u.postgres.CreateTeacher(ctx, &teacher)
		if err != nil {
			return Output{}, err
		}
	}

	return Output{ID: user.ID}, nil
}
