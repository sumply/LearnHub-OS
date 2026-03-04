package create_group

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type Group interface {
	repository.Saver[domain.Group]
	AddTeacher(ctx context.Context, groupID uuid.UUID, teacherID uuid.UUID) error
	AddStudents(ctx context.Context, groupID uuid.UUID, studentIDs uuid.UUIDs) error
}

type User interface {
	repository.Geter[domain.User]
	repository.Lister[domain.User]
}

type UseCase struct {
	group Group
	user  User
}

func New(group Group, user User) *UseCase {
	return &UseCase{
		group: group,
		user:  user,
	}
}

func (u *UseCase) CreateGroup(ctx context.Context, identity usecase.Identity, input *Input) (Output, error) {
	if identity.Role() != domain.RoleAdmin {
		return Output{}, usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}

	group := domain.Group{
		ID:   uuid.New(),
		Name: input.Name,
	}

	err := u.group.Save(ctx, group)
	if err != nil {
		return Output{}, err
	}

	err = u.handleOptionFields(ctx, &group, input)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: group.ID}, nil
}

func (u *UseCase) handleOptionFields(ctx context.Context, group *domain.Group, input *Input) error {
	if input.CuratorID != uuid.Nil {
		user, err := u.user.Get(ctx, input.CuratorID)
		if err != nil {
			return err
		}
		if user.Role != domain.RoleTeacher {
			return usecase.NewValidationError(fmt.Errorf("curator is not teacher"))
		}
		err = u.group.AddTeacher(ctx, group.ID, user.ID)
		if err != nil {
			return err
		}
	}

	if len(input.StudentIDs) != 0 {
		users, err := u.user.List(ctx, input.StudentIDs)
		if err != nil {
			return err
		}

		studentIDs := make(uuid.UUIDs, 0, len(users))
		for _, user := range users {
			if user.Role != domain.RoleStudent {
				return usecase.NewValidationError(fmt.Errorf("user is not student"))
			}
			studentIDs = append(studentIDs, user.ID)
		}

		err = u.group.AddStudents(ctx, group.ID, studentIDs)
		if err != nil {
			return err
		}
	}

	return nil
}
