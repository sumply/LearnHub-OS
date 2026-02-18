package create_user

import (
	"context"
	"server/internal/domain"
	"server/internal/usecase"

	"github.com/google/uuid"
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

func (u *UseCase) CreateUser(ctx context.Context, input Input) (Output, error) {
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
		if err := u.createStudent(ctx, input, user); err != nil {
			return Output{}, err
		}
	case domain.RoleTeacher:
		if err := u.createTeacher(ctx, input, user); err != nil {
			return Output{}, err
		}
	}

	return Output{ID: user.ID}, nil
}

func (u *UseCase) createStudent(ctx context.Context, input Input, user domain.User) error {
	if input.GroupID == uuid.Nil {
		err := u.postgres.CreateUser(ctx, &user)
		if err != nil {
			return err
		}
	} else {
		student, err := domain.NewStudent(user, input.GroupID)
		if err != nil {
			return err
		}

		err = u.postgres.CreateStudent(ctx, &student)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UseCase) createTeacher(ctx context.Context, input Input, user domain.User) error {
	if len(input.GroupIDs) == 0 && len(input.SubjectIDs) == 0 {
		err := u.postgres.CreateUser(ctx, &user)
		if err != nil {
			return err
		}
	} else {
		teacher, err := domain.NewTeacher(user, input.GroupIDs)
		if err != nil {
			return err
		}

		err = u.postgres.CreateTeacher(ctx, &teacher)
		if err != nil {
			return err
		}
	}
	return nil
}
