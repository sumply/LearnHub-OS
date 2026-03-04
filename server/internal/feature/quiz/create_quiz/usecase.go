package create_quiz

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"
	"server/internal/pkg/usecase"
)

type TeacherRepository interface {
	repository.Geter[domain.Teacher]
}

type QuizRepository interface {
	repository.Saver[domain.Quiz]
}

type UseCase struct {
	teacher TeacherRepository
	quiz    QuizRepository
}

func New() *UseCase {
	return &UseCase{}
}

func (u *UseCase) CreateQuiz(ctx context.Context, identity usecase.Identity, input *Input) (Output, error) {
	if identity.Role() != domain.RoleAdmin || identity.Role() != domain.RoleTeacher {
		return Output{}, usecase.NewAuthError(
			fmt.Errorf("user is not admin or teacher"),
		)
	}

	quiz, err := u.createQuiz(input)
	if err != nil {
		return Output{}, err
	}

	err = u.quiz.Save(ctx, quiz)
	if err != nil {
		return Output{}, err
	}

	return Output{ID: quiz.ID}, nil
}

func (u *UseCase) validateAccess(_ context.Context, identity usecase.Identity, _ *Input) error {
	switch identity.Role() {
	case domain.RoleAdmin:
		return nil
	case domain.RoleTeacher:
		return nil
	default:
		return usecase.NewAuthError(
			fmt.Errorf("user has not access to creating quiz"),
		)
	}
}

func (u *UseCase) createQuiz(input *Input) (domain.Quiz, error) {
	questions := make([]domain.Question, len(input.Questions))
	for i := range questions {
		question, err := u.createQuestion(&input.Questions[i])
		if err != nil {
			return domain.Quiz{}, err
		}

		questions[i] = question
	}
	quiz, err := domain.NewQuiz(
		input.OwnerID,
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

func (u *UseCase) createQuestion(input *dto.Question) (domain.Question, error) {
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
