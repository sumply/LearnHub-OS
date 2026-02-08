package create_user

import (
	"context"
	"fmt"
	"server/internal/domain"
	"time"

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

func (u *UseCase) CreateUser(ctx context.Context, input *Input) (Output, error) {
	user := domain.User{
		ID:        uuid.New(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      domain.UserRole(input.Role),
		Email:     input.Email,
		PwdHash:   "hash",
		CreatedAt: time.Now().UTC(),
	}

	switch domain.UserRole(input.Role) {

	case domain.RoleAdmin:
		err := u.postgres.CreateUser(ctx, &user)
		if err != nil {
			return Output{}, err
		}

	case domain.RoleStudent:
		if input.GroupID == nil {
			return Output{}, fmt.Errorf("groupID is nil")
		}
		student := domain.Student{
			User:  user,
			Group: *input.GroupID,
		}
		err := u.postgres.CreateStudent(ctx, &student)
		if err != nil {
			return Output{}, err
		}

	case domain.RoleTeacher:
		teacher := domain.Teacher{
			User:     user,
			Subjects: input.SubjectIDs,
			Groups:   input.GroupIDs,
		}
		err := u.postgres.CreateTeacher(ctx, &teacher)
		if err != nil {
			return Output{}, err
		}
	}

	return Output{ID: user.ID}, nil
}
