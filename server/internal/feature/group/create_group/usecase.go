package create_group

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type DomainRepository interface {
	AddGroup(context.Context, *domain.Group) error
	GetUser(context.Context, uuid.UUID) (domain.User, error)
	ListUserByIDs(context.Context, uuid.UUIDs) ([]domain.User, error)
	AddTeacherToGroup(ctx context.Context, groupID uuid.UUID, teacherID uuid.UUID) error
	AddStudentsToGroup(ctx context.Context, groupID uuid.UUID, studentIDs uuid.UUIDs) error
}

type UseCase struct {
	dRepository DomainRepository
}

func New(repostiory DomainRepository) *UseCase {
	return &UseCase{
		dRepository: repostiory,
	}
}

func (u *UseCase) CreateGroup(ctx context.Context, input *Input) (Output, error) {
	group := domain.Group{
		ID:   uuid.New(),
		Name: input.Name,
	}

	err := u.dRepository.AddGroup(ctx, &group)
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
	if input.CuratorID != nil {
		user, err := u.dRepository.GetUser(ctx, *input.CuratorID)
		if err != nil {
			return err
		}
		if user.Role != domain.RoleTeacher {
			return usecase.NewValidationError(fmt.Errorf("curator is not teacher"))
		}
		err = u.dRepository.AddTeacherToGroup(ctx, group.ID, user.ID)
		if err != nil {
			return err
		}
	}

	if len(input.StudentIDs) != 0 {
		users, err := u.dRepository.ListUserByIDs(ctx, input.StudentIDs)
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

		err = u.dRepository.AddStudentsToGroup(ctx, group.ID, studentIDs)
		if err != nil {
			return err
		}
	}

	return nil
}
