package create_user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/usecase"

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

func (uc *UseCase) CreateUser(ctx context.Context, identity usecase.Identity, input Input) (Output, error) {
	if identity.Role() != domain.RoleAdmin {
		return Output{}, usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

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

	err = uc.saveUserByRole(ctx, input, user)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: user.ID}, nil
}

func (uc *UseCase) saveUserByRole(ctx context.Context, input Input, user domain.User) error {
	switch domain.UserRole(input.Role) {

	case domain.RoleAdmin:
		err := uc.postgres.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

	case domain.RoleStudent:
		if err := uc.createStudent(ctx, input, user); err != nil {
			return err
		}
	case domain.RoleTeacher:
		if err := uc.createTeacher(ctx, input, user); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) createStudent(ctx context.Context, input Input, user domain.User) error {
	if input.GroupID == uuid.Nil {
		err := uc.postgres.CreateUser(ctx, &user)
		if err != nil {
			return err
		}
	} else {
		student, err := domain.NewStudent(user, input.GroupID)
		if err != nil {
			return err
		}

		err = uc.postgres.CreateStudent(ctx, &student)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) createTeacher(ctx context.Context, input Input, user domain.User) error {
	if len(input.GroupIDs) == 0 && len(input.SubjectIDs) == 0 {
		err := uc.postgres.CreateUser(ctx, &user)
		if err != nil {
			return err
		}
	} else {
		teacher, err := domain.NewTeacher(user, input.GroupIDs)
		if err != nil {
			return err
		}

		err = uc.postgres.CreateTeacher(ctx, &teacher)
		if err != nil {
			return err
		}
	}
	return nil
}
