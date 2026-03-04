package create_quiz

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"

	"github.com/google/uuid"
)

type Teacher interface {
	repository.Geter[domain.Teacher]
}

type Quiz interface {
	repository.Saver[domain.Quiz]
}

type UseCase struct {
	teacher Teacher
	quiz    Quiz
}

func New(t Teacher, q Quiz) *UseCase {
	return &UseCase{
		teacher: t,
		quiz:    q,
	}
}

func (u *UseCase) CreateQuiz(ctx context.Context, identity usecase.Identity, input Input) (Output, error) {
	err := u.validateAccess(ctx, identity, input)
	if err != nil {
		return Output{}, err
	}

	quiz, err := u.createQuiz(identity.ID(), input)
	if err != nil {
		return Output{}, err
	}

	err = u.quiz.Save(ctx, quiz)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: quiz.ID}, nil
}

func (u *UseCase) validateAccess(ctx context.Context, identity usecase.Identity, input Input) error {
	switch identity.Role() {
	case domain.RoleAdmin:
		return nil
	case domain.RoleTeacher:
		teacher, err := u.teacher.Get(ctx, identity.ID())
		if err != nil {
			return err
		}
		err = teacher.CheckGroupsAllowed(input.GroupIDs)
		if err != nil {
			return usecase.NewAuthError(err)
		}
		return nil
	default:
		return usecase.NewAuthError(
			fmt.Errorf("user has not access to creating quiz"),
		)
	}
}

func (u *UseCase) createQuiz(ownerID uuid.UUID, input Input) (domain.Quiz, error) {
	questions := make([]domain.Question, len(input.Questions))
	for i := range questions {
		question, err := u.createQuestion(input.Questions[i])
		if err != nil {
			return domain.Quiz{}, err
		}

		questions[i] = question
	}
	quiz, err := domain.NewQuiz(
		ownerID,
		input.SubjectID,
		input.Title,
		input.Summary,
		questions,
		input.GroupIDs,
		input.MaxAttempts,
		input.Deadline,
	)
	if err != nil {
		return domain.Quiz{}, err
	}

	return quiz, nil
}

func (u *UseCase) createQuestion(input dto.Question) (domain.Question, error) {
	question, err := domain.NewQuestion(
		input.Text,
		input.Details.Domain,
		input.Score,
	)
	if err != nil {
		return domain.Question{}, err
	}

	return question, nil
}
