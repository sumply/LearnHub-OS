package create_user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Transaction interface {
	Execute(context.Context, func(context.Context) error) error
}

type Postgres interface {
	CreateUser(context.Context, *domain.User) error
	CreateStudent(context.Context, *domain.Student) error
	BindTeacherToGroups(context.Context, uuid.UUID, uuid.UUIDs) error
	BindTeacherToSubjects(context.Context, uuid.UUID, uuid.UUIDs) error
}

type UseCase struct {
	postgres    Postgres
	transaction Transaction
}

func New(postgres Postgres, transaction Transaction) *UseCase {
	return &UseCase{
		postgres:    postgres,
		transaction: transaction,
	}
}

func (u *UseCase) CreateUser(ctx context.Context, input *Input) (Output, error) {
	role, err := domain.NewUserRole(input.Role)
	if err != nil {
		return Output{}, err
	}

	var userID uuid.UUID

	err = u.transaction.Execute(ctx, func(ctx context.Context) error {
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
			return err
		}

		switch role {
		case domain.RoleStudent:
			if input.GroupID != nil {
				student := domain.Student{
					User:  user.ID,
					Group: *input.GroupID,
				}
				err := u.postgres.CreateStudent(ctx, &student)
				if err != nil {
					return err
				}
			}
		case domain.RoleTeacher:
			if input.GroupIDs == nil || input.SubjectIDs == nil {
				return fmt.Errorf("groupIDs or subjectIDs: %w", domain.ErrInvalid)
			}
			err := u.postgres.BindTeacherToGroups(ctx, user.ID, input.GroupIDs)
			if err != nil {
				return err
			}
			err = u.postgres.BindTeacherToSubjects(ctx, user.ID, input.SubjectIDs)
			if err != nil {
				return err
			}
		}

		userID = user.ID

		return nil
	})
	if err != nil {
		return Output{}, err
	}

	return Output{ID: userID}, nil
}
