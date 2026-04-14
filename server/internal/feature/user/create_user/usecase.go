package create_user

import (
	"context"
	"fmt"
	"log/slog"
	"server/internal/domain"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"
	"server/pkg/logger"

	"github.com/google/uuid"
)

type User interface {
	repository.Saver[domain.User]
}

type Teacher interface {
	repository.Saver[domain.Teacher]
}

type Student interface {
	repository.Saver[domain.Student]
}

type UseCase struct {
	user    User
	teacher Teacher
	student Student
}

func New(
	user User,
	teacher Teacher,
	student Student) *UseCase {
	return &UseCase{
		user:    user,
		teacher: teacher,
		student: student,
	}
}

func (uc *UseCase) CreateUser(ctx context.Context, identity usecase.Identity, input Input) (Output, error) {
	log := logger.FromCtx(ctx)
	log = log.With(
		usecase.IdentityToSlogAttr(identity),
		slog.Group("input",
			slog.String("first_name", input.FirstName),
			slog.String("last_name", input.LastName),
			slog.String("email", input.Email),
			slog.String("role", input.Role),
			slog.String("group_id", input.GroupID.String()),
			slog.Any("group_ids", input.GroupIDs),
			slog.Any("subject_id", input.SubjectIDs.Strings()),
		),
	)
	if err := uc.validateAccess(identity); err != nil {
		log.WarnContext(ctx, "Failed validating access", slog.String("error", err.Error()))
		return Output{}, err
	}

	user, err := uc.makeUser(input)
	if err != nil {
		log.WarnContext(ctx, "Failed making user", slog.String("error", err.Error()))
		return Output{}, err
	}

	if err := uc.save(ctx, user, input); err != nil {
		return Output{}, err
	}

	return Output{ID: user.ID}, nil
}

func (uc *UseCase) validateAccess(identity usecase.Identity) error {
	if identity.Role() != domain.RoleAdmin {
		return usecase.NewAuthError(
			fmt.Errorf("user is not admin"),
		)
	}
	return nil
}

func (uc *UseCase) makeUser(input Input) (domain.User, error) {
	user, err := domain.NewUser(
		input.FirstName,
		input.LastName,
		input.Email,
		"hash",
		domain.UserRole(input.Role),
	)
	if err != nil {
		return domain.User{}, usecase.NewValidationError(err)
	}
	return user, nil
}

func (uc *UseCase) save(ctx context.Context, user domain.User, input Input) error {
	switch user.Role {
	case domain.RoleAdmin:
		return uc.saveAsUser(ctx, user)
	case domain.RoleTeacher:
		return uc.saveAsTeacher(ctx, user, input.SubjectIDs)
	case domain.RoleStudent:
		return uc.saveAsStudent(ctx, user, input.GroupID)
	}
	return nil
}

func (uc *UseCase) saveAsUser(ctx context.Context, user domain.User) error {
	return uc.user.Save(ctx, user)
}

func (uc *UseCase) saveAsTeacher(ctx context.Context, user domain.User, groupIDs []uuid.UUID) error {
	if len(groupIDs) == 0 {
		return uc.saveAsUser(ctx, user)
	}
	teacher, err := domain.NewTeacher(user, groupIDs)
	if err != nil {
		return usecase.NewValidationError(err)
	}
	return uc.teacher.Save(ctx, teacher)
}

func (uc *UseCase) saveAsStudent(ctx context.Context, user domain.User, groupID uuid.UUID) error {
	if groupID == uuid.Nil {
		return uc.saveAsUser(ctx, user)
	}
	student, err := domain.NewStudent(user, groupID)
	if err != nil {
		return usecase.NewValidationError(err)
	}
	return uc.student.Save(ctx, student)
}
